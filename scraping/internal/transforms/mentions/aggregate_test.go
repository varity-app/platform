package mentions

import (
	"testing"

	"gotest.tools/assert"

	"github.com/varity-app/platform/scraping/internal/data/bigquery"
)

const (
	sourceDummy    = "source1"
	subredditDummy = "subreddit1"
)

// TestNoMentions tests the aggregate function with no mentions
func TestNoMentions(t *testing.T) {
	memberships := []bigquery.CohortMembership{}
	mentions := []AugmentedMention{}

	aggs := aggregate(memberships, mentions)

	check := map[[4]string][3]int{}

	assert.DeepEqual(t, aggs, check)
}

// TestNoMemberships is a unit test
func TestNoMemberships(t *testing.T) {
	memberships := []bigquery.CohortMembership{}
	mentions := []AugmentedMention{
		{Symbol: "AAPL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: false},
		{Symbol: "AAPL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true},
		{Symbol: "GOOGL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true},
	}

	aggs := aggregate(memberships, mentions)

	check := map[[4]string][3]int{
		{"AAPL", subredditDummy, cohortAll, sourceDummy}:  {2, 2, 1},
		{"GOOGL", subredditDummy, cohortAll, sourceDummy}: {1, 1, 1},
	}

	assert.DeepEqual(t, aggs, check)
}

// TestOneMembership is a unit test
func TestOneMembership(t *testing.T) {
	memberships := []bigquery.CohortMembership{
		{AuthorID: "author1", Subreddit: subredditDummy, Cohort: "cohort1"},
	}
	mentions := []AugmentedMention{
		{Symbol: "AAPL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: false, AuthorID: "author1"},
		{Symbol: "GOOGL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true},
		{Symbol: "GOOGL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true},
	}

	aggs := aggregate(memberships, mentions)

	check := map[[4]string][3]int{
		{"AAPL", subredditDummy, "cohort1", sourceDummy}:  {1, 1, 0},
		{"AAPL", subredditDummy, cohortAll, sourceDummy}:  {1, 1, 0},
		{"GOOGL", subredditDummy, cohortAll, sourceDummy}: {2, 2, 2},
	}

	assert.DeepEqual(t, aggs, check)
}

// TestThreeMembership is a unit test
func TestThreeMembership(t *testing.T) {
	memberships := []bigquery.CohortMembership{
		{AuthorID: "author1", Subreddit: subredditDummy, Cohort: "cohort1"},
		{AuthorID: "author1", Subreddit: subredditDummy, Cohort: "cohort2"},
		{AuthorID: "author1", Subreddit: subredditDummy, Cohort: "cohort3"},
	}
	mentions := []AugmentedMention{
		{Symbol: "AAPL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: false, AuthorID: "author1"},
		{Symbol: "AAPL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true, AuthorID: "author2"},
		{Symbol: "GOOGL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true, AuthorID: "author3"},
		{Symbol: "GOOGL", Subreddit: subredditDummy, Source: sourceDummy, Targeted: true, Inquisitive: true, AuthorID: "author3"},
	}

	aggs := aggregate(memberships, mentions)

	check := map[[4]string][3]int{
		{"AAPL", subredditDummy, "cohort1", sourceDummy}:  {1, 1, 0},
		{"AAPL", subredditDummy, "cohort2", sourceDummy}:  {1, 1, 0},
		{"AAPL", subredditDummy, "cohort3", sourceDummy}:  {1, 1, 0},
		{"AAPL", subredditDummy, cohortAll, sourceDummy}:  {2, 2, 1},
		{"GOOGL", subredditDummy, cohortAll, sourceDummy}: {2, 2, 2},
	}

	assert.DeepEqual(t, aggs, check)
}
