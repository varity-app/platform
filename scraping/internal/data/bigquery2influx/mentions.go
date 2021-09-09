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
	// MentionsTableReddit is the name of the bigquery table containing Reddit ticker mentions
	MentionsTableReddit string = "reddit_mentions"
)

// Mention represents a redditor's membership to a cohort.
type Mention struct {
	Symbol      string    `json:"symbol" bigquery:"symbol"`
	Timestamp   time.Time `json:"timestamp" bigquery:"timestamp"`
	Targeted    bool      `json:"targeted" bigquery:"targeted"`
	Inquisitive bool      `json:"inquisitive" bigquery:"inquisitive"`
	AuthorID    string    `json:"author_id" bigquery:"author_id"`
	Subreddit   string    `json:"subreddit" bigquery:"subreddit"`
	Source      string    `json:"source" bigquery:"source"`
}

// MentionsRepoOpts is a helper object used to pass configuration parameters
// to NewMentionsRepo()
type MentionsRepoOpts struct {
}

// MentionsRepo is an abstraction layer for retreiving Mention objects from BigQuery.
type MentionsRepo struct {
	bqClient *bigquery.Client
	rdb      *redis.Client
}

// NewMentionsRepo is a constructor that returns a new CohortMebershipRepo
func NewMentionsRepo(bqClient *bigquery.Client, redisOpts *redis.Options) *MentionsRepo {
	return &MentionsRepo{
		bqClient: bqClient,
		rdb:      redis.NewClient(redisOpts),
	}
}

// GetMentions fetches all mentions for a specific cohort table for a specific month.
// tableName should be a full table reference of the format `project.dataset.tablename`
func (r *MentionsRepo) GetMentions(ctx context.Context, tableName string, year, month, day, hour int) ([]Mention, error) {

	// Check Redis cache
	redisKey := fmt.Sprintf("%s.%d.%d.%d", tableName, year, month, hour)
	mentions, err := r.checkCache(ctx, redisKey)
	if err == nil { // Return value retrieved from cache
		return mentions, nil
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
			symbol, timestamp, targeted, inquisitive, author_id, subreddit, source
		FROM %s
		WHERE
			TIMESTAMP_TRUNC(timestamp, HOUR) = @hour;
	`, tableName))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "hour", Value: fmt.Sprintf("%d-%d-%d %d:00:00 UTC", year, month, day, hour)},
	}

	// Run interactive query
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipRepo.Query: %v", err)
	}

	// Iterate over query results
	mentions = []Mention{}
	for {
		var mention Mention
		err = it.Next(&mention)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("cohortMembershipRepo.QueryResults: %v", err)
		}

		mentions = append(mentions, mention)
	}

	// Save query results to redis
	err = r.saveCache(ctx, redisKey, mentions)
	if err != nil {
		return nil, fmt.Errorf("mentionsRepo.SaveCache: %v", err)
	}

	return mentions, nil
}

// checkCache queries the Redis cache to see if a queries results have been stored.
func (r *MentionsRepo) checkCache(ctx context.Context, key string) ([]Mention, error) {

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var mentions []Mention
	err = json.Unmarshal([]byte(val), &mentions)
	if err != nil {
		return nil, fmt.Errorf("cohortMentions.UnmarshalJSON: %v", err)
	}

	return mentions, nil
}

// saveCache saves a query's results to the Redis cache.
func (r *MentionsRepo) saveCache(ctx context.Context, key string, mentions []Mention) error {

	result, err := json.Marshal(mentions)
	if err != nil {
		return fmt.Errorf("mentions.MarshalJSON: %v", err)
	}

	err = r.rdb.Set(ctx, key, result, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
