package bigquery2influx

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"

	"github.com/go-redis/redis/v8"
)

const (
	// CohortTableRedditorsAge is the name of the bigquery table containing the Redditors Age cohort
	CohortTableRedditorsAge string = "redditors_age_cohort_memberships"

	// CohortTableRedditorsPopularity is the name of the bigquery table containing the Redditors Popularity cohort
	CohortTableRedditorsPopularity string = "redditors_popularity_cohort_memberships"

	// CohortTableRedditorsFrequencies is the name of the bigquery table containing the Redditors Posting Frequencies cohort
	CohortTableRedditorsFrequencies string = "redditors_frequencies_cohort_memberships"
)

// CohortMembership represents a redditor's membership to a cohort.
type CohortMembership struct {
	AuthorID  string    `json:"author_id" bigquery:"author_id"`
	Subreddit string    `json:"subreddit" bigquery:"subreddit"`
	Month     time.Time `json:"month" bigquery:"month"`
	Cohort    string    `json:"cohort" bigquery:"cohort"`
}

// CohortMembershipRepoOpts is a helper object used to pass configuration parameters
// to NewCohortMembershipRepo()
type CohortMembershipRepoOpts struct {
}

// CohortMembershipRepo is an abstraction layer for retreiving CohortMebership objects from BigQuery.
type CohortMembershipRepo struct {
	bqClient *bigquery.Client
	rdb      *redis.Client
}

// NewCohortMembershipRepo is a constructor that returns a new CohortMebershipRepo
func NewCohortMembershipRepo(bqClient *bigquery.Client, redisOpts *redis.Options) *CohortMembershipRepo {
	return &CohortMembershipRepo{
		bqClient: bqClient,
		rdb:      redis.NewClient(redisOpts),
	}
}

// GetMemberships fetches all memberships for a specific cohort table for a specific month.
// tableName should be a full table reference of the format `project.dataset.tablename`
func (r *CohortMembershipRepo) GetMemberships(ctx context.Context, tableName string, year, month int) ([]CohortMembership, error) {

	// Check Redis cache
	redisKey := fmt.Sprintf("%s_%d_%d", tableName, year, month)
	memberships, err := r.checkCache(ctx, redisKey)
	if err == nil { // Return value tretrieved from cache
		return memberships, nil
	} else if err != redis.Nil {
		return nil, fmt.Errorf("cohortMembershipRepo.CheckCache: %v", err)
	}

	// Check that tableName is of valid format
	matched, err := regexp.Match(`\w+[.]\w+[.]\w+`, []byte(tableName))
	if err != nil {
		return nil, fmt.Errorf("regexp.Match: %v", err)
	}
	if !matched {
		return nil, fmt.Errorf("cohortMembership.CheckTableName: invalid tableName: %s", tableName)
	}

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			author_id, subreddit, month, cohort
		FROM %s
		WHERE
			month = @month;
	`, tableName))
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
	err = r.saveCache(ctx, redisKey, memberships)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.SaveCache: %v", err)
	}

	return memberships, nil
}

// checkCache queries the Redis cache to see if a query's results have been stored.
func (r *CohortMembershipRepo) checkCache(ctx context.Context, key string) ([]CohortMembership, error) {

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var memberships []CohortMembership
	err = json.Unmarshal([]byte(val), &memberships)
	if err != nil {
		return nil, fmt.Errorf("cohortMemberships.UnmarshalJSON: %v", err)
	}

	return memberships, nil
}

// saveCache saves a query's results to the Redis cache.
func (r *CohortMembershipRepo) saveCache(ctx context.Context, key string, memberships []CohortMembership) error {

	result, err := json.Marshal(memberships)
	if err != nil {
		return fmt.Errorf("cohortMemberships.MarshalJSON: %v", err)
	}

	err = r.rdb.Set(ctx, key, result, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
