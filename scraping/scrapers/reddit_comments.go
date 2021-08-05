package scrapers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//RedditCommentsScraper Scrapes comments from reddit
type RedditCommentsScraper struct {
	redditClient *reddit.Client
	memory       *Memory
}

//NewRedditCommentsScraper Initialize a new RedditCommentScraper
func NewRedditCommentsScraper(redditCredentials reddit.Credentials, memory *Memory) (*RedditCommentsScraper, error) {
	redditClient, err := reddit.NewClient(redditCredentials, reddit.WithUserAgent(RedditUserAgent))
	if err != nil {
		return nil, fmt.Errorf("reddit.NewClient: %v", err)
	}

	return &RedditCommentsScraper{
		redditClient: redditClient,
		memory:       memory,
	}, nil
}

// Close disconnects the scraper connection
func (scraper *RedditCommentsScraper) Close() error {
	return scraper.memory.Close()
}

// Scrape scrapes new submissions from reddit
func (scraper *RedditCommentsScraper) Scrape(ctx context.Context, subreddit string, limit int) ([]*rpb.RedditComment, error) {
	comments, _, err := scraper.fetchComments(ctx, subreddit, &reddit.ListOptions{
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("redditCommentScraper.Fetch: %v", err)
	}

	protoComments := scraper.convert(comments)

	// Get list of comment ids to check with memory
	var ids []string
	for _, comment := range protoComments {
		ids = append(ids, comment.GetCommentId())
	}

	// Check memory for which comments have been seen
	unseenIdxs, err := scraper.memory.CheckNewItems(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("redditCommentScraper.CheckMemory: %v", err)
	}

	// Get list of unseen comments with the list of indices produced by memory
	var unseenComments []*rpb.RedditComment
	for _, idx := range unseenIdxs {
		unseenComments = append(unseenComments, protoComments[idx])
	}

	return unseenComments, nil
}

// CommitSeen commits comments to memory
func (scraper *RedditCommentsScraper) CommitSeen(ctx context.Context, comments []*rpb.RedditComment) error {
	// Get list of comments ids to check with memory
	var ids []string
	for _, comments := range comments {
		ids = append(ids, comments.GetCommentId())
	}

	err := scraper.memory.SaveItems(ctx, ids)
	if err != nil {
		return fmt.Errorf("redditCommentsScraper.CommitMemory: %v", err)
	}

	return nil
}

func (scraper *RedditCommentsScraper) fetchComments(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Comment, *reddit.Response, error) {
	// Add options to HTTP path
	path, err := addOptions("r/"+subreddit+"/comments", opts)
	if err != nil {
		return nil, nil, err
	}

	// Create the request
	req, err := scraper.redditClient.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Run the request
	t := new(thing)
	resp, err := scraper.redditClient.Do(ctx, req, t)
	if err != nil {
		return nil, nil, err
	}

	// Parse data to comments
	l, _ := t.Listing()

	return l.Comments(), resp, nil
}

func (scraper *RedditCommentsScraper) convert(comments []*reddit.Comment) []*rpb.RedditComment {
	var protoComments []*rpb.RedditComment

	for _, comment := range comments {
		protoComment := &rpb.RedditComment{
			CommentId:    comment.ID,
			SubmissionId: comment.ParentID,
			Subreddit:    comment.SubredditName,
			Timestamp:    timestamppb.New(comment.Created.Time),
			Body:         comment.Body,
			Author:       comment.Author,
			AuthorId:     comment.AuthorID,
			Permalink:    &comment.Permalink,
		}

		protoComments = append(protoComments, protoComment)
	}

	return protoComments
}
