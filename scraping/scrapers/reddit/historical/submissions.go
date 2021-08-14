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

type pushshiftSubmission struct {
	SubmissionID string `json:"id"`
	Subreddit    string `json:"subreddit"`
	Title        string `json:"title"`
	CreatedUTC   int64  `json:"created_utc"`
	Body         string `json:"selftext"`
	Author       string `json:"author"`
	AuthorID     string `json:"author_fullname"`
	IsSelf       bool   `json:"is_self"`
	Permalink    string `json:"permalink"`
	URL          string `json:"url"`
}

type pushshiftSubmissionResponse struct {
	Submissions []pushshiftSubmission `json:"data"`
}

// SubmissionsScraper scrapes historical reddit submissions from the PushshiftAPI
type SubmissionsScraper struct {
	client *http.Client
	memory *scrapers.Memory
}

// NewSubmissionsScraper initializes a new historical submissions scraper
func NewSubmissionsScraper(memory *scrapers.Memory) (*SubmissionsScraper, error) {

	return &SubmissionsScraper{
		client: &http.Client{},
		memory: memory,
	}, nil
}

// Scrape scrapes reddit submissions from the Pushshift API
func (scraper *SubmissionsScraper) Scrape(ctx context.Context, subreddit string, before, after time.Time, limit int) ([]*rpb.RedditSubmission, error) {

	// Create request
	req, err := http.NewRequest("GET", PushshiftBaseURL+"/reddit/submission/search", nil)
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
		return nil, fmt.Errorf("submissionsScraper.httpRequest: %v", err)
	}

	// Extract submissions from body
	submissions, err := scraper.extract(resp)
	if err != nil {
		return nil, err
	}

	// Filter seen submissions
	unseenSubmissions, err := scraper.filter(ctx, submissions)
	if err != nil {
		return nil, err
	}

	return unseenSubmissions, nil
}

// CommitSeen commits seen submissions to memory
func (scraper *SubmissionsScraper) CommitSeen(ctx context.Context, submissions []*rpb.RedditSubmission) error {
	// Get list of submission ids to check with memory
	var ids []string
	for _, submission := range submissions {
		ids = append(ids, submission.GetSubmissionId())
	}

	err := scraper.memory.SaveItems(ctx, ids)
	if err != nil {
		return fmt.Errorf("submissionsScraper.CommitMemory: %v", err)
	}

	return nil
}

// Filter out seen submissions
func (scraper *SubmissionsScraper) filter(ctx context.Context, submissions []*rpb.RedditSubmission) ([]*rpb.RedditSubmission, error) {

	// Get list of submission ids to check with memory
	var ids []string
	for _, submission := range submissions {
		ids = append(ids, submission.GetSubmissionId())
	}

	// Check memory for which submissions have been seen
	unseenIdxs, err := scraper.memory.CheckNewItems(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("submissionsScraper.CheckMemory: %v", err)
	}

	// Get list of unseen submissions with the list of indices produced by memory
	var unseenSubmissions []*rpb.RedditSubmission
	for _, idx := range unseenIdxs {
		unseenSubmissions = append(unseenSubmissions, submissions[idx])
	}

	return unseenSubmissions, nil
}

// Extract reddit submissions from the response body
func (scraper *SubmissionsScraper) extract(resp *http.Response) ([]*rpb.RedditSubmission, error) {
	defer resp.Body.Close()

	// Read respose body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("scraper.ReadResponseBody: %v", err)
	}

	// Unmarshal response
	parsedResponse := pushshiftSubmissionResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("scraper.Unmarshal: %v", err)
	}

	// Convert response to protobuf
	submissions := scraper.convert(parsedResponse.Submissions)

	return submissions, nil
}

// Convert reddit submissions to protobuf messages
func (scraper *SubmissionsScraper) convert(submissions []pushshiftSubmission) []*rpb.RedditSubmission {
	var protoSubmissions []*rpb.RedditSubmission

	for _, submission := range submissions {
		timestamp := time.Unix(submission.CreatedUTC, 0)

		protoSubmission := &rpb.RedditSubmission{
			SubmissionId: submission.SubmissionID,
			Subreddit:    submission.Subreddit,
			Title:        submission.Title,
			Timestamp:    timestamppb.New(timestamp),
			Body:         submission.Body,
			Author:       submission.Author,
			AuthorId:     submission.AuthorID,
			IsSelf:       &submission.IsSelf,
			Permalink:    &submission.Permalink,
			Url:          &submission.URL,
		}
		protoSubmissions = append(protoSubmissions, protoSubmission)
	}

	return protoSubmissions
}
