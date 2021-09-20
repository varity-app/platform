package live

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	rpb "github.com/varity-app/platform/scraping/api/reddit/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/scrapers"
)

// RedditSubmissionsScraper Scrapes comments from reddit
type RedditSubmissionsScraper struct {
	redditClient *reddit.Client
	memory       *scrapers.Memory
}

// NewRedditSubmissionsScraper Initialize a new RedditSubmissionScraper
func NewRedditSubmissionsScraper(redditCredentials reddit.Credentials, memory *scrapers.Memory) (*RedditSubmissionsScraper, error) {
	redditClient, err := reddit.NewClient(redditCredentials, reddit.WithUserAgent(RedditUserAgent))
	if err != nil {
		return nil, fmt.Errorf("reddit.NewClient: %v", err)
	}

	return &RedditSubmissionsScraper{
		redditClient: redditClient,
		memory:       memory,
	}, nil
}

// Scrape Scrape new submissions from reddit
func (scraper *RedditSubmissionsScraper) Scrape(ctx context.Context, subreddit string, limit int) ([]*rpb.RedditSubmission, error) {
	submissions, _, err := scraper.redditClient.Subreddit.NewPosts(ctx, subreddit, &reddit.ListOptions{
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("redditSubmissionsScraper.Fetch: %v", err)
	}

	protoSubmissions := scraper.convert(submissions)

	// Get list of submission ids to check with memory
	var ids []string
	for _, submission := range protoSubmissions {
		ids = append(ids, submission.GetSubmissionId())
	}

	// Check memory for which submissions have been seen
	unseenIdxs, err := scraper.memory.CheckNewItems(ctx, common.RedditSubmissions, ids)
	if err != nil {
		return nil, fmt.Errorf("redditSubmissionsScraper.CheckMemory: %v", err)
	}

	// Get list of unseen submissions with the list of indices produced by memory
	var unseenSubmissions []*rpb.RedditSubmission
	for _, idx := range unseenIdxs {
		unseenSubmissions = append(unseenSubmissions, protoSubmissions[idx])
	}

	return unseenSubmissions, nil
}

// CommitSeen Commits submissions to memory
func (scraper *RedditSubmissionsScraper) CommitSeen(ctx context.Context, submissions []*rpb.RedditSubmission) error {
	// Get list of submission ids to check with memory
	var ids []string
	for _, submission := range submissions {
		ids = append(ids, submission.GetSubmissionId())
	}

	err := scraper.memory.SaveItems(ctx, common.RedditSubmissions, ids)
	if err != nil {
		return fmt.Errorf("redditSubmissionsScraper.CommitMemory: %v", err)
	}

	return nil
}

// Convert reddit submissions to protobuf messages
func (scraper *RedditSubmissionsScraper) convert(submissions []*reddit.Post) []*rpb.RedditSubmission {
	var protoSubmissions []*rpb.RedditSubmission

	for _, submission := range submissions {
		protoSubmission := &rpb.RedditSubmission{
			SubmissionId: submission.ID,
			Subreddit:    submission.SubredditName,
			Title:        submission.Title,
			Timestamp:    timestamppb.New(submission.Created.Time),
			Body:         submission.Body,
			Author:       submission.Author,
			AuthorId:     submission.AuthorID,
			IsSelf:       &submission.IsSelfPost,
			Permalink:    &submission.Permalink,
			Url:          &submission.URL,
		}
		protoSubmissions = append(protoSubmissions, protoSubmission)
	}

	return protoSubmissions
}
