package transforms

import (
	"github.com/VarityPlatform/scraping/common"
	pb "github.com/VarityPlatform/scraping/protobuf/common"
	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"
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
