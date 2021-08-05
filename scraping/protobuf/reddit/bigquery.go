package reddit

import (
	"cloud.google.com/go/bigquery"
)

// Save implements the ValueSaver interface.  This is used for bigquery.
func (submission *RedditSubmission) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"submission_id": submission.GetSubmissionId(),
		"subreddit":     submission.GetSubreddit(),
		"title":         submission.GetTitle(),
		"timestamp":     submission.GetTimestamp().AsTime(),
		"body":          submission.GetBody(),
		"author":        submission.GetAuthor(),
		"author_id":     submission.GetAuthorId(),
		"is_self":       submission.GetIsSelf(),
		"permalink":     submission.GetPermalink(),
		"url":           submission.GetUrl(),
	}, bigquery.NoDedupeID, nil
}

func (comment *RedditComment) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"comment_id":    comment.GetCommentId(),
		"submission_id": comment.GetSubmissionId(),
		"subreddit":     comment.GetSubreddit(),
		"author":        comment.GetAuthor(),
		"author_id":     comment.GetAuthorId(),
		"timestamp":     comment.GetTimestamp().AsTime(),
		"body":          comment.GetBody(),
		"permalink":     comment.GetPermalink(),
	}, bigquery.NoDedupeID, nil
}
