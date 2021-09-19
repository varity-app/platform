package mentions

import (
	"time"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data/bigquery"
)

// Threshold at which a post is classified as targeting a specific ticker.
// A post is targeted if a single ticker is mentioned X times more than any other ticker
const targetedThreshold int = 3

// Threshold at which a post is classified as inquisitive.  The metris is calculated by
// number of question marks / number of words in post
const inquisitiveThreshold float32 = 0.1

const sourceRedditComment = "reddit-comment"
const sourceRedditSubmission = "reddit-submission"

// AugmentedMention represents an augmented ticker mention.
type AugmentedMention struct {
	Symbol      string
	Timestamp   time.Time
	Targeted    bool
	Inquisitive bool
	AuthorID    string
	Subreddit   string
	Source      string
}

// AugmentMentions aggregates a given hour's ticker mentions, reddit submissions, and comments
// and generates the Targeted and Inquisitive annotations for each mention.
func AugmentMentions(
	mentions []bigquery.Mention,
	submissions []bigquery.RedditSubmission,
	comments []bigquery.RedditComment,
) []AugmentedMention {
	var augmented []AugmentedMention

	commentMentions, submissionMentions := groupMentions(mentions)

	augmented = append(augmented, augmentComments(commentMentions, comments)...)
	augmented = append(augmented, augmentSubmissions(submissionMentions, submissions)...)

	return augmented
}

// Group mentions by their parent posts
func groupMentions(mentions []bigquery.Mention) (map[string][]bigquery.Mention, map[string][]bigquery.Mention) {

	// Generate a map of post id -> mentions
	commentMentions := map[string][]bigquery.Mention{}
	submissionMentions := map[string][]bigquery.Mention{}
	for _, mention := range mentions {

		isComment := mention.ParentSource == common.ParentSourceRedditComment
		isSubmission := mention.ParentSource == common.ParentSourceRedditSubmissionTitle
		isSubmission = isSubmission || mention.ParentSource == common.ParentSourceRedditSubmissionBody

		// Assign mention to correct map
		if isComment {
			commentMentions[mention.ParentID] = append(commentMentions[mention.ParentID], mention)
		} else if isSubmission {
			submissionMentions[mention.ParentID] = append(submissionMentions[mention.ParentID], mention)
		}
	}

	return commentMentions, submissionMentions
}

// Augment comment mentions
func augmentComments(commentMentions map[string][]bigquery.Mention, comments []bigquery.RedditComment) []AugmentedMention {
	var augmented []AugmentedMention

	for _, comment := range comments {
		mentions := commentMentions[comment.CommentID]

		// Ensure that there are mention in the map for this comment
		if len(mentions) == 0 {
			continue
		}

		// Generate frequencies for each ticker
		tickerCounts := map[string]int{}
		totalCounts := 0
		for _, mention := range mentions {
			tickerCounts[mention.Symbol] += mention.SymbolCounts + mention.ShortNameCounts
			totalCounts += mention.SymbolCounts + mention.ShortNameCounts
		}

		// Calculate the most common
		mostFrequent := getMostFrequent(tickerCounts)
		otherMentions := totalCounts - tickerCounts[mostFrequent]

		// Determine if comment is targeting a specific ticker
		targeted := false
		if len(tickerCounts) == 1 {
			targeted = true
		} else if otherMentions != 0 { // Check if one ticker was mentioned X times more than all other tickers combined.
			targeted = tickerCounts[mostFrequent]/otherMentions >= targetedThreshold
		}

		// Determine if a comment is inquisitive
		wordCount := mentions[0].WordCount
		questionCount := mentions[0].QuestionMarkCount

		inquisitive := float32(questionCount)/float32(wordCount) >= inquisitiveThreshold

		// Append new augmented mentions
		for symbol := range tickerCounts {
			augmented = append(augmented, AugmentedMention{
				Symbol:      symbol,
				Timestamp:   comment.Timestamp,
				Targeted:    targeted && symbol == mostFrequent,
				Inquisitive: inquisitive,
				AuthorID:    comment.AuthorID,
				Subreddit:   comment.Subreddit,
				Source:      sourceRedditComment,
			})
		}

	}

	return augmented
}

// Augment submission mentions
func augmentSubmissions(submissionMentions map[string][]bigquery.Mention, submissions []bigquery.RedditSubmission) []AugmentedMention {
	var augmented []AugmentedMention

	for _, submission := range submissions {
		mentions := submissionMentions[submission.SubmissionID]

		// Ensure that there are mention in the map for this submission
		if len(mentions) == 0 {
			continue
		}

		// Set word counts individually for title and body
		titleWordCount := 0
		titleQuestionCount := 0
		bodyWordCount := 0
		bodyQuestionCount := 0

		// Generate frequencies for each ticker
		tickerCounts := map[string]int{}
		totalCounts := 0
		titleSymbolsCount := 0
		for _, mention := range mentions {
			tickerCounts[mention.Symbol] += mention.SymbolCounts + mention.ShortNameCounts
			totalCounts += mention.SymbolCounts + mention.ShortNameCounts

			// Increment number of symbols found in post title
			if mention.ParentSource == common.ParentSourceRedditSubmissionTitle {
				titleSymbolsCount++
				titleWordCount = mention.WordCount
				titleQuestionCount = mention.QuestionMarkCount
			} else {
				bodyWordCount = mention.WordCount
				bodyQuestionCount = mention.QuestionMarkCount
			}
		}

		// Calculate the most common
		mostFrequent := getMostFrequent(tickerCounts)
		otherMentions := totalCounts - tickerCounts[mostFrequent]

		// Determine if submission is targeting a specific ticker
		targeted := false
		if titleSymbolsCount == 1 || len(tickerCounts) == 1 { // Only one symbol mentioned in title, or only one in entire post
			targeted = true
		} else if otherMentions > 0 { // Check if one ticker was mentioned X times more than all other tickers combined.
			targeted = tickerCounts[mostFrequent]/otherMentions >= targetedThreshold
		}

		// Determine if a submission is inquisitive
		wordCount := titleWordCount + bodyWordCount
		questionCount := titleQuestionCount + bodyQuestionCount

		inquisitive := titleWordCount > 0 && float32(titleQuestionCount)/float32(titleWordCount) >= inquisitiveThreshold
		inquisitive = inquisitive || float32(questionCount)/float32(wordCount) >= inquisitiveThreshold

		// Append new augmented mentions
		for symbol := range tickerCounts {
			augmented = append(augmented, AugmentedMention{
				Symbol:      symbol,
				Timestamp:   submission.Timestamp,
				Targeted:    targeted && symbol == mostFrequent,
				Inquisitive: inquisitive,
				AuthorID:    submission.AuthorID,
				Subreddit:   submission.Subreddit,
				Source:      sourceRedditSubmission,
			})
		}

	}

	return augmented
}

// Find the most frequently mentioned ticker from a map of frequencies
func getMostFrequent(frequencies map[string]int) string {
	mostFrequent := ""

	for ticker, counts := range frequencies {
		if mostFrequent == "" {
			mostFrequent = ticker
		} else {
			if frequencies[mostFrequent] < counts {
				mostFrequent = ticker
			}
		}
	}

	return mostFrequent
}
