package bigquery2influx

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
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

// MentionsRepo is an abstraction layer for retreiving Mention objects from BigQuery.
type MentionsRepo struct {
	bqClient      *bigquery.Client
	bigqueryTable string
}

// NewMentionsRepo is a constructor that returns a new CohortMebershipRepo
func NewMentionsRepo(bqClient *bigquery.Client, deploymentMode string) *MentionsRepo {
	return &MentionsRepo{
		bqClient:      bqClient,
		bigqueryTable: wrapDBTTable(deploymentMode, MentionsTableReddit),
	}
}

// GetMentions fetches all mentions for a specific cohort table for a specific month.
// tableName should be a full table reference of the format `project.dataset.tablename`
func (r *MentionsRepo) GetMentions(ctx context.Context, year, month, day, hour int) ([]Mention, error) {

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			symbol, timestamp, targeted, inquisitive, author_id, subreddit, source
		FROM %s
		WHERE
			TIMESTAMP_TRUNC(timestamp, HOUR) = @hour;
	`, r.bigqueryTable))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "hour", Value: fmt.Sprintf("%d-%d-%d %d:00:00 UTC", year, month, day, hour)},
	}

	// Run interactive query
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("mentionsRepo.Query: %v", err)
	}

	// Iterate over query results
	mentions := []Mention{}
	for {
		var mention Mention
		err = it.Next(&mention)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("mentionsRepo.QueryResults: %v", err)
		}

		mentions = append(mentions, mention)
	}

	return mentions, nil
}
