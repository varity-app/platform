package bigquery

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
)

const (
	// MentionsTableReddit is the name of the bigquery table containing Reddit ticker mentions
	MentionsTableReddit string = "ticker_mentions_v2"
)

// Mention represents a redditor's membership to a cohort.
type Mention struct {
	Symbol            string    `json:"symbol" bigquery:"symbol"`
	ParentSource      string    `json:"parent_source" bigquery:"parent_source"`
	ParentID          string    `json:"parent_id" bigquery:"parent_id"`
	Timestamp         time.Time `json:"timestamp" bigquery:"timestamp"`
	SymbolCounts      int       `json:"symbol_counts" bigquery:"symbol_counts"`
	ShortNameCounts   int       `json:"short_name_counts" bigquery:"short_name_counts"`
	WordCount         int       `json:"word_count" bigquery:"word_count"`
	QuestionMarkCount int       `json:"question_mark_count" bigquery:"question_mark_count"`
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
		bigqueryTable: wrapTable(deploymentMode, MentionsTableReddit),
	}
}

// Get fetches all mentions for a specific month.
func (r *MentionsRepo) Get(ctx context.Context, year, month, day, hour int) ([]Mention, error) {

	// Formulate query
	query := r.bqClient.Query(fmt.Sprintf(`
		SELECT
			symbol, parent_source, parent_id, timestamp, symbol_counts, short_name_counts, word_count, question_mark_count
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
