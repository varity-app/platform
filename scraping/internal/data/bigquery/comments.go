package bigquery

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
)

const (
	// TableComments is the name of the bigquery table containing Reddit comments
	TableComments string = "reddit_comments_v2"
)

// RedditComment represents a redditor's membership to a cohort.
type RedditComment struct {
	CommentID string    `json:"symbol" bigquery:"comment_id"`
	Subreddit string    `json:"subreddit" bigquery:"subreddit"`
	Timestamp time.Time `json:"timestamp" bigquery:"timestamp"`
	AuthorID  string    `json:"author_id" bigquery:"author_id"`
}

// RedditCommentsRepo is an abstraction layer for retreiving RedditComment objects from BigQuery.
type RedditCommentsRepo struct {
	bqClient      *bigquery.Client
	bigqueryTable string
}

// NewRedditCommentsRepo is a constructor that returns a new CohortMebershipRepo
func NewRedditCommentsRepo(bqClient *bigquery.Client, deploymentMode string) *RedditCommentsRepo {
	return &RedditCommentsRepo{
		bqClient:      bqClient,
		bigqueryTable: wrapTable(deploymentMode, TableComments),
	}
}

// Get fetches all comments for a specific hour in time.
func (r *RedditCommentsRepo) Get(ctx context.Context, year, month, day, hour int) ([]RedditComment, error) {

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			comment_id, subreddit, timestamp, author_id
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
		return nil, fmt.Errorf("commentsRepo.Query: %v", err)
	}

	// Iterate over query results
	comments := []RedditComment{}
	for {
		var comment RedditComment
		err = it.Next(&comment)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("commentsRepo.QueryResults: %v", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
