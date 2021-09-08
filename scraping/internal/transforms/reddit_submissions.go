package transforms

import (
	rpb "github.com/varity-app/platform/scraping/api/reddit/v1"
	pb "github.com/varity-app/platform/scraping/api/ticker_mentions/v1"
	"github.com/varity-app/platform/scraping/internal/common"
)

// TransformRedditSubmission extracts ticker mentions from a reddit submission
func TransformRedditSubmission(extractor *TickerExtractor, submission *rpb.RedditSubmission) []pb.TickerMention {
	// Parse tickers
	titleMentions := extractor.ExtractTickerMentions(submission.Title, submission.SubmissionId, common.ParentSourceRedditSubmissionTitle, submission.Timestamp)
	bodyMentions := extractor.ExtractTickerMentions(submission.GetBody(), submission.SubmissionId, common.ParentSourceRedditSubmissionBody, submission.Timestamp) // Concat tickers into one array

	// Concat mentions into one array
	allMentions := append(titleMentions, bodyMentions...)

	return allMentions
}
