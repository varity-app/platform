package transforms

import (
	"github.com/VarityPlatform/scraping/common"
	pb "github.com/VarityPlatform/scraping/protobuf/common"
	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"
)

// TransformRedditComment extracts ticker mentions from a reddit comment
func TransformRedditComment(extractor *TickerExtractor, comment *rpb.RedditComment) []pb.TickerMention {
	// Parse tickers
	mentions := extractor.ExtractTickerMentions(comment.Body, comment.CommentId, common.ParentSourceRedditComment, comment.Timestamp)

	return mentions
}
