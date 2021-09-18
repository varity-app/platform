package bigquery2influx

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"google.golang.org/api/iterator"
	"gorm.io/gorm"

	"cloud.google.com/go/bigquery"
)

const (
	// cohortTableRedditorsAge is the name of the bigquery table containing the Redditors Age cohort
	cohortTableRedditorsAge string = "redditors_age_cohort_memberships"

	// cohortTableRedditorsPopularity is the name of the bigquery table containing the Redditors Popularity cohort
	cohortTableRedditorsPopularity string = "redditors_popularity_cohort_memberships"

	// cohortTableRedditorsFrequencies is the name of the bigquery table containing the Redditors Posting Frequencies cohort
	cohortTableRedditorsFrequencies string = "redditors_frequencies_cohort_memberships"

	// Batch size to use when inserting cohort memberships into Postgres
	cohortsBatchSize = 1000
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
	pdb            *gorm.DB
	deploymentMode string
}

// NewCohortMembershipRepo is a constructor that returns a new CohortMebershipRepo
func NewCohortMembershipRepo(bqClient *bigquery.Client, pdb *gorm.DB, deploymentMode string) *CohortMembershipRepo {
	return &CohortMembershipRepo{
		bqClient:       bqClient,
		pdb:            pdb,
		deploymentMode: deploymentMode,
	}
}

// GetMemberships fetches all memberships for a specific cohort table for a specific month.
// tableName should be a full table reference of the format `project.dataset.tablename`
func (r *CohortMembershipRepo) GetMemberships(ctx context.Context, year, month int) ([]CohortMembership, error) {

	var allMemberships []CohortMembership

	// Query each individual table
	for table, _ := range cohortsTableMap {
		memberships, err := r.queryCohort(ctx, table, year, month)
		if err != nil {
			return nil, err
		}

		allMemberships = append(allMemberships, memberships...)
	}

	return allMemberships, nil
}

func (r *CohortMembershipRepo) queryCohort(ctx context.Context, table string, year, month int) ([]CohortMembership, error) {
	cohortPrefix := cohortsTableMap[table]
	memberships, err := r.checkCache(ctx, cohortPrefix, year, month)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.CheckCache: %v", err)
	}

	if len(memberships) != 0 {
		return memberships, nil
	}

	log.Println("querying bigquery...")

	// Check that table is of valid format
	matched, err := regexp.Match(`\w+`, []byte(table))
	if err != nil {
		return nil, fmt.Errorf("regexp.Match: %v", err)
	}
	if !matched {
		return nil, fmt.Errorf("cohortMembership.CheckTableName: invalid table name: %s", table)
	}

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			author_id, subreddit, month, cohort
		FROM %s
		WHERE
			month = @month;
	`, wrapDBTTable(r.deploymentMode, table)))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "month", Value: fmt.Sprintf("%d-%d-01", year, month)},
	}

	// Run interactive query
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.Query: %v", err)
	}

	// Iterate over query results
	memberships = []CohortMembership{}
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

	// Save query results to redis
	err = r.saveCache(memberships)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.SaveCache: %v", err)
	}

	return memberships, nil
}

// checkCache queries the Firestore cache to see if a queries results have been stored.
func (r *CohortMembershipRepo) checkCache(ctx context.Context, cohortPrefix string, year, month int) ([]CohortMembership, error) {

	// Create time.Time representation of month
	ts := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Query postgres
	var memberships []CohortMembership
	result := r.pdb.Where("month = ? AND cohort LIKE ?", ts, cohortPrefix+"%").Find(&memberships)
	if result.Error != nil {
		return nil, fmt.Errorf("postgres.Get: %v", result.Error)
	}

	return memberships, nil
}

// saveCache saves a query's results to the Firestore cache.
func (r *CohortMembershipRepo) saveCache(memberships []CohortMembership) error {
	result := r.pdb.CreateInBatches(memberships, cohortsBatchSize)

	return result.Error
}
