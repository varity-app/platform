package historical

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/VarityPlatform/scraping/scrapers"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"
)

type pushshiftComment struct {
	CommentID    string `json:"id"`
	SubmissionID string `json:"parent_id"`
	Subreddit    string `json:"subreddit"`
	CreatedUTC   int64  `json:"created_utc"`
	Body         string `json:"body"`
	Author       string `json:"author"`
	AuthorID     string `json:"author_fullname"`
	Permalink    string `json:"permalink"`
}

type pushshiftCommentResponse struct {
	Comments []pushshiftComment `json:"data"`
}

// CommentsScraper scrapes historical reddit comments from the PushshiftAPI
type CommentsScraper struct {
	client *http.Client
	memory *scrapers.Memory
}

// NewCommentsScraper initializes a new historical comments scraper
func NewCommentsScraper(opts scrapers.MemoryOpts, memory *scrapers.Memory) (*CommentsScraper, error) {
	return &CommentsScraper{
		client: &http.Client{},
		memory: memory,
	}, nil
}

// Scrape scrapes reddit comments from the Pushshift API
func (scraper *CommentsScraper) Scrape(ctx context.Context, subreddit string, before, after time.Time, limit int) ([]*rpb.RedditComment, error) {

	// Create request
	req, err := http.NewRequest("GET", PushshiftBaseURL+"/reddit/comment/search", nil)
	if err != nil {
		return nil, err
	}

	// Add token parameter
	q := req.URL.Query()
	q.Add("subreddit", subreddit)
	q.Add("before", fmt.Sprint(before.Unix()))
	q.Add("after", fmt.Sprint(after.Unix()))
	q.Add("size", fmt.Sprint(limit))
	req.URL.RawQuery = q.Encode()

	resp, err := scraper.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("commentsScraper.httpRequest: %v", err)
	}

	// Extract comments from body
	comments, err := scraper.extract(resp)
	if err != nil {
		return nil, err
	}

	// Filter seen comments
	unseenComments, err := scraper.filter(ctx, comments)
	if err != nil {
		return nil, err
	}

	return unseenComments, nil
}

// CommitSeen commits seen comments to memory
func (scraper *CommentsScraper) CommitSeen(ctx context.Context, comments []*rpb.RedditComment) error {
	// Get list of comment ids to check with memory
	var ids []string
	for _, comment := range comments {
		ids = append(ids, comment.GetCommentId())
	}

	err := scraper.memory.SaveItems(ctx, ids)
	if err != nil {
		return fmt.Errorf("commentsScraper.CommitMemory: %v", err)
	}

	return nil
}

// Filter out seen comments
func (scraper *CommentsScraper) filter(ctx context.Context, comments []*rpb.RedditComment) ([]*rpb.RedditComment, error) {

	// Get list of comment ids to check with memory
	var ids []string
	for _, comment := range comments {
		ids = append(ids, comment.GetCommentId())
	}

	// Check memory for which comments have been seen
	unseenIdxs, err := scraper.memory.CheckNewItems(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("commentsScraper.CheckMemory: %v", err)
	}

	// Get list of unseen comments with the list of indices produced by memory
	var unseenComments []*rpb.RedditComment
	for _, idx := range unseenIdxs {
		unseenComments = append(unseenComments, comments[idx])
	}

	return unseenComments, nil
}

// Extract reddit comments from the response body
func (scraper *CommentsScraper) extract(resp *http.Response) ([]*rpb.RedditComment, error) {
	defer resp.Body.Close()

	// Read respose body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("scraper.ReadResponseBody: %v", err)
	}

	// Unmarshal response
	parsedResponse := pushshiftCommentResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("scraper.Unmarshal: %v", err)
	}

	// Convert response to protobuf
	comments := scraper.convert(parsedResponse.Comments)

	return comments, nil
}

// Convert reddit comments to protobuf messages
func (scraper *CommentsScraper) convert(comments []pushshiftComment) []*rpb.RedditComment {
	var protoComments []*rpb.RedditComment

	for _, comment := range comments {
		timestamp := time.Unix(comment.CreatedUTC, 0)

		protoComment := &rpb.RedditComment{
			CommentId:    comment.CommentID,
			SubmissionId: comment.SubmissionID,
			Subreddit:    comment.Subreddit,
			Timestamp:    timestamppb.New(timestamp),
			Body:         comment.Body,
			Author:       comment.Author,
			AuthorId:     comment.AuthorID,
			Permalink:    &comment.Permalink,
		}
		protoComments = append(protoComments, protoComment)
	}

	return protoComments
}
