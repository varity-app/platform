package bigquery

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
)

const (
	// TableSubmissions is the name of the bigquery table containing Reddit submissions
	TableSubmissions string = "reddit_submissions_v2"
)

// RedditSubmission represents a redditor's membership to a cohort.
type RedditSubmission struct {
	SubmissionID string    `json:"symbol" bigquery:"submission_id"`
	Subreddit    string    `json:"subreddit" bigquery:"subreddit"`
	Timestamp    time.Time `json:"timestamp" bigquery:"timestamp"`
	AuthorID     string    `json:"author_id" bigquery:"author_id"`
}

// RedditSubmissionsRepo is an abstraction layer for retreiving RedditSubmission objects from BigQuery.
type RedditSubmissionsRepo struct {
	bqClient      *bigquery.Client
	bigqueryTable string
}

// NewRedditSubmissionsRepo is a constructor that returns a new CohortMebershipRepo
func NewRedditSubmissionsRepo(bqClient *bigquery.Client, deploymentMode string) *RedditSubmissionsRepo {
	return &RedditSubmissionsRepo{
		bqClient:      bqClient,
		bigqueryTable: wrapTable(deploymentMode, TableSubmissions),
	}
}

// Get fetches all submissions for a specific hour in time.
func (r *RedditSubmissionsRepo) Get(ctx context.Context, year, month, day, hour int) ([]RedditSubmission, error) {

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			submission_id, subreddit, timestamp, author_id
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
		return nil, fmt.Errorf("submissionsRepo.Query: %v", err)
	}

	// Iterate over query results
	submissions := []RedditSubmission{}
	for {
		var submission RedditSubmission
		err = it.Next(&submission)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("submissionsRepo.QueryResults: %v", err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}
