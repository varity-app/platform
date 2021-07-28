package main

import (
	"regexp"

	"github.com/VarityPlatform/scraping/common"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonPB "github.com/VarityPlatform/scraping/protobuf/common"
	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
)

var questionRegex *regexp.Regexp = regexp.MustCompile(`\?`)

// Process reddit submissions
func procRedditSubmission(submission *redditPB.RedditSubmission, allTickers []common.IEXTicker) []commonPB.TickerMention {
	// Parse tickers
	titleMentions := procPost(submission.Title, submission.SubmissionId, common.PARENT_SOURCE_REDDIT_SUBMISSION_TITLE, submission.Timestamp, allTickers)
	bodyMentions := procPost(submission.Body, submission.SubmissionId, common.PARENT_SOURCE_REDDIT_SUBMISSION_BODY, submission.Timestamp, allTickers)

	// Concat tickers into one array
	allMentions := append(titleMentions, bodyMentions...)

	return allMentions
}

// Process a post text field
func procPost(s string, id string, parentSource string, timestamp *timestamppb.Timestamp, allTickers []common.IEXTicker) []commonPB.TickerMention {
	if s == "" {
		return []commonPB.TickerMention{}
	}

	// Extract tickers and shortname mentions from string
	tickers := extractTickersString(s, allTickers)
	nameTickers := extractShortNamesString(s, tickers)

	// Calculate frequencies
	uniqTickers, tickerFrequencies := calcTickerFrequency(tickers)
	_, nameFrequencies := calcTickerFrequency(nameTickers)

	// Calculate extra metrics
	questionCount := len(questionRegex.FindAllString(s, -1))
	wordCount := len(wordRegex.FindAllString(s, -1))

	// Generate list of ticker mentions
	mentions := []commonPB.TickerMention{}
	for _, ticker := range uniqTickers {
		mentions = append(mentions, commonPB.TickerMention{
			Symbol:            ticker.Symbol,
			ParentId:          id,
			ParentSource:      parentSource,
			Timestamp:         timestamp,
			SymbolCounts:      uint32(tickerFrequencies[ticker.Symbol]),
			ShortNameCounts:   uint32(nameFrequencies[ticker.Symbol]),
			WordCount:         uint32(wordCount),
			QuestionMarkCount: uint32(questionCount),
		})
	}

	return mentions
}
