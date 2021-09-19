package bigquery

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/logging"
	"google.golang.org/api/iterator"

	lru "github.com/hashicorp/golang-lru"

	"cloud.google.com/go/bigquery"
)

const (
	// cohortTableRedditorsAge is the name of the bigquery table containing the Redditors Age cohort
	cohortTableRedditorsAge string = "redditors_age_cohort_memberships"

	// cohortTableRedditorsPopularity is the name of the bigquery table containing the Redditors Popularity cohort
	cohortTableRedditorsPopularity string = "redditors_popularity_cohort_memberships"

	// cohortTableRedditorsFrequencies is the name of the bigquery table containing the Redditors Posting Frequencies cohort
	cohortTableRedditorsFrequencies string = "redditors_frequencies_cohort_memberships"
)

// Map a table to it's corresponding cohort prefix.
// E.g. redditors_age_cohort_memberships -> "age*"
var cohortsTableMap map[string]string = map[string]string{
	cohortTableRedditorsAge:         "age",
	cohortTableRedditorsPopularity:  "popularity",
	cohortTableRedditorsFrequencies: "post_frequency",
}

// CohortMembership represents a redditor's membership to a cohort.
type CohortMembership struct {
	AuthorID  string    `json:"author_id" bigquery:"author_id"`
	Subreddit string    `json:"subreddit" bigquery:"subreddit"`
	Month     time.Time `json:"month" bigquery:"month"`
	Cohort    string    `json:"cohort" bigquery:"cohort"`
}

// CohortMembershipRepo is an abstraction layer for retreiving CohortMebership objects from BigQuery.
type CohortMembershipRepo struct {
	bqClient       *bigquery.Client
	lru            *lru.Cache
	deploymentMode string
	logger         *logging.Logger
}

// NewCohortMembershipRepo is a constructor that returns a new CohortMebershipRepo
func NewCohortMembershipRepo(bqClient *bigquery.Client, logger *logging.Logger, deploymentMode string) (*CohortMembershipRepo, error) {

	// Create LRU cache
	lruCache, err := lru.New(12)
	if err != nil {
		return nil, fmt.Errorf("lru.New: %v", err)
	}

	return &CohortMembershipRepo{
		bqClient:       bqClient,
		lru:            lruCache,
		deploymentMode: deploymentMode,
		logger:         logger,
	}, nil
}

// Get fetches all cohort memberships for a specific month.
func (r *CohortMembershipRepo) Get(ctx context.Context, year, month int) ([]CohortMembership, error) {

	var allMemberships []CohortMembership

	// Query each individual table
	for table := range cohortsTableMap {
		memberships, err := r.queryCohort(ctx, table, year, month)
		if err != nil {
			return nil, err
		}
		r.logger.Debug(fmt.Sprintf("Found %d cohort memberships in table %s", len(memberships), table))

		allMemberships = append(allMemberships, memberships...)
	}

	return allMemberships, nil
}

func (r *CohortMembershipRepo) queryCohort(ctx context.Context, table string, year, month int) ([]CohortMembership, error) {
	cohortPrefix := cohortsTableMap[table]

	// Check cach table for memberships
	memberships, err := r.checkCache(ctx, cohortPrefix, year, month)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.CheckCache: %v", err)
	}

	if len(memberships) != 0 {
		return memberships, nil
	}

	// Check that table is of valid format
	matched, err := regexp.Match(`\w+`, []byte(table))
	if err != nil {
		return nil, fmt.Errorf("regexp.Match: %v", err)
	}
	if !matched {
		return nil, fmt.Errorf("cohortMembership.CheckTableName: invalid table name: %s", table)
	}

	log.Println("Running cache query...")

	// Run caching query, since no memberships where found
	err = r.saveCache(ctx, table, year, month)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.SaveCache: %v", err)
	}

	log.Println("Done running cache query...")

	// Check cach table again for memberships
	memberships, err = r.checkCache(ctx, cohortPrefix, year, month)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.CheckCache: %v", err)
	}

	return memberships, nil
}

// checkCache queries the Firestore cache to see if a queries results have been stored.
func (r *CohortMembershipRepo) checkCache(ctx context.Context, cohortPrefix string, year, month int) ([]CohortMembership, error) {

	// Check LRU cache for query results
	lruKey := fmt.Sprintf("%s:%d:%d", cohortPrefix, year, month)
	cachedResult, ok := r.lru.Get(lruKey)
	if ok && cachedResult != nil {
		memberships, ok := cachedResult.([]CohortMembership)
		if !ok {
			return nil, fmt.Errorf("lru.ConvertResult: unable to convert cache to []CohortMembership")
		}
		return memberships, nil
	}
	r.logger.Debug("Memberships not found in LRU cache.  Querying bigquery...")

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			author_id, subreddit, month, cohort
		FROM %s
		WHERE
			month = @month AND cohort LIKE @like;
	`, wrapTable(r.deploymentMode, common.BigqueryTableMemberships)))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "month", Value: fmt.Sprintf("%d-%d-01", year, month)},
		{Name: "like", Value: cohortPrefix + "%"},
	}

	// Run interactive query
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.Query: %v", err)
	}

	// Iterate over query results
	var memberships []CohortMembership
	for {
		var membership CohortMembership
		err = it.Next(&membership)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("cohortMembershipRepo.QueryResults: %v", err)
		}

		memberships = append(memberships, membership)
	}

	// Save query result to LRU
	r.lru.Add(lruKey, memberships)

	return memberships, nil
}

// saveCache runs the BigQuery query and saves results in cache table.
func (r *CohortMembershipRepo) saveCache(
	ctx context.Context,
	sourceTable string,
	year, month int,
) error {

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(
		`
		INSERT %s (author_id, subreddit, month, cohort)
		SELECT
			author_id, subreddit, month, cohort
		FROM %s
		WHERE
			month = @month;
		`,
		wrapTable(r.deploymentMode, common.BigqueryTableMemberships),
		wrapDBTTable(r.deploymentMode, sourceTable),
	))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "month", Value: fmt.Sprintf("%d-%d-01", year, month)},
	}

	// Run query
	_, err := query.Read(ctx)
	if err != nil {
		return fmt.Errorf("cohortMembershipRepo.Query: %v", err)
	}

	return nil
}
