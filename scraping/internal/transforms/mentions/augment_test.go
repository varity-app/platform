package mentions

import (
	"reflect"
	"testing"
	"time"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data/bigquery"
)

var now time.Time = time.Now()

// TestCommentsEmpty is a unit test
func TestCommentsEmpty(t *testing.T) {
	mentions := []bigquery.Mention{}
	comments := []bigquery.RedditComment{}

	commentMentions, _ := groupMentions(mentions)
	augmented := augmentComments(commentMentions, comments)

	var check []AugmentedMention

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestCommentsOneTargeted is a unit test
func TestCommentsOneTargeted(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "GOOGL",
			WordCount:         15,
			QuestionMarkCount: 1,
		},
	}
	comments := []bigquery.RedditComment{
		{CommentID: "id1", Timestamp: now},
	}

	commentMentions, _ := groupMentions(mentions)
	augmented := augmentComments(commentMentions, comments)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    true,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestCommentsOneTargeted is a unit test
func TestCommentsOneTargetedInquisitive(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "GOOGL",
			WordCount:         10,
			QuestionMarkCount: 1,
		},
	}
	comments := []bigquery.RedditComment{
		{CommentID: "id1", Timestamp: now},
	}

	commentMentions, _ := groupMentions(mentions)
	augmented := augmentComments(commentMentions, comments)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: true,
			Targeted:    true,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestCommentsTwoNonTargeted is a unit test
func TestCommentsTwoNonTargeted(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
	}
	comments := []bigquery.RedditComment{
		{CommentID: "id1", Timestamp: now},
	}

	commentMentions, _ := groupMentions(mentions)
	augmented := augmentComments(commentMentions, comments)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
		{
			Symbol:      "AAPL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestCommentsThreeTargeted is a unit test
func TestCommentsThreeTargeted(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "GOOGL",
			SymbolCounts:      5,
			ShortNameCounts:   2,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditComment,
			Symbol:            "GME",
			SymbolCounts:      1,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
	}
	comments := []bigquery.RedditComment{
		{CommentID: "id1", Timestamp: now},
	}

	commentMentions, _ := groupMentions(mentions)
	augmented := augmentComments(commentMentions, comments)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    true,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
		{
			Symbol:      "AAPL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
		{
			Symbol:      "GME",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditComment,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsEmpty is a unit test
func TestSubmissionsEmpty(t *testing.T) {
	mentions := []bigquery.Mention{}
	submissions := []bigquery.RedditSubmission{}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	var check []AugmentedMention

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsOne is a unit test
func TestSubmissionsOne(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         15,
			QuestionMarkCount: 1,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         30,
			QuestionMarkCount: 2,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    true,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsInquisitiveTitle is a unit test
func TestSubmissionsInquisitiveTitle(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         15,
			QuestionMarkCount: 2,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         30,
			QuestionMarkCount: 2,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: true,
			Targeted:    true,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsInquisitive is a unit test
func TestSubmissionsInquisitiveBody(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         20,
			QuestionMarkCount: 3,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: true,
			Targeted:    true,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsTwoNonTargeted is a unit test
func TestSubmissionsTwoNonTargeted(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
		{
			Symbol:      "AAPL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsTwoTargeted is a unit test
func TestSubmissionsTwoTargeted(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			ShortNameCounts:   1,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    true,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
		{
			Symbol:      "AAPL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}

// TestSubmissionsTwoTargetedAgain is a unit test
func TestSubmissionsTwoTargetedAgain(t *testing.T) {
	mentions := []bigquery.Mention{
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "GOOGL",
			SymbolCounts:      2,
			ShortNameCounts:   1,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "GOOGL",
			SymbolCounts:      1,
			ShortNameCounts:   2,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionTitle,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			ShortNameCounts:   0,
			WordCount:         10,
			QuestionMarkCount: 0,
		},
		{
			ParentID:          "id1",
			ParentSource:      common.ParentSourceRedditSubmissionBody,
			Symbol:            "AAPL",
			SymbolCounts:      1,
			ShortNameCounts:   0,
			WordCount:         20,
			QuestionMarkCount: 0,
		},
	}

	submissions := []bigquery.RedditSubmission{
		{SubmissionID: "id1", Timestamp: now},
	}

	_, submissionsMentions := groupMentions(mentions)
	augmented := augmentSubmissions(submissionsMentions, submissions)

	check := []AugmentedMention{
		{
			Symbol:      "GOOGL",
			Inquisitive: false,
			Targeted:    true,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
		{
			Symbol:      "AAPL",
			Inquisitive: false,
			Targeted:    false,
			Source:      sourceRedditSubmission,
			Timestamp:   now,
		},
	}

	if !reflect.DeepEqual(augmented, check) {
		t.Errorf("Assertion failed: got %v, want: %v", augmented, check)
	}
}
