package transforms

import (
	rpb "github.com/varity-app/platform/scraping/api/reddit/v1"
	pb "github.com/varity-app/platform/scraping/api/ticker_mentions/v1"
	"github.com/varity-app/platform/scraping/internal/common"
)

// TransformRedditComment extracts ticker mentions from a reddit comment
func TransformRedditComment(extractor *TickerExtractor, comment *rpb.RedditComment) []pb.TickerMention {
	// Parse tickers
	mentions := extractor.ExtractTickerMentions(comment.Body, comment.CommentId, common.ParentSourceRedditComment, comment.Timestamp)

	return mentions
}
