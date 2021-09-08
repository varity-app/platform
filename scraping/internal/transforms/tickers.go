package transforms

import (
	"regexp"

	pb "github.com/varity-app/platform/scraping/api/ticker_mentions/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TickerSpamThreshold is describes the threshold of number of ticker mentions past which
// posts are flagged as spam and no tickers are returned.
const TickerSpamThreshold int = 5

var urls *regexp.Regexp = regexp.MustCompile(`https?:\/\/.*[\r\n]*`)
var alphaNumeric *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 \n\.]`)
var capsSpam *regexp.Regexp = regexp.MustCompile(`([A-Z]{1,3}.?[A-Z]{1,3})(\W[A-Z].?[A-Z]+)+`)
var tickerRegex *regexp.Regexp = regexp.MustCompile(`[A-Z][A-Z0-9.]*[A-Z0-9]`)
var wordRegex *regexp.Regexp = regexp.MustCompile(`\w+`)
var questionRegex *regexp.Regexp = regexp.MustCompile(`\?`)

// TickerExtractor is an object that extracts tickers from a string, given
// a dictionary of tickers.
type TickerExtractor struct {
	tickerList []common.IEXTicker
}

// NewTickerExtractor initializes a new TickerExtractor
func NewTickerExtractor(tickerList []common.IEXTicker) *TickerExtractor {
	return &TickerExtractor{
		tickerList: tickerList,
	}
}

// ExtractTickers extracts tickers from a string
func (extractor *TickerExtractor) ExtractTickers(s string) []common.IEXTicker {

	// Remove urls
	s = urls.ReplaceAllString(s, "")

	// Remove alphanumeric characters
	s = alphaNumeric.ReplaceAllString(s, "")

	// Remove caps spam
	s = capsSpam.ReplaceAllString(s, "")

	// Find ticker look-a-likes
	tickers := tickerRegex.FindAllString(s, -1)

	// Cross reference with provided list of valid tickers
	validTickers := []common.IEXTicker{}
	for _, ticker := range tickers {
		blacklisted := false
		// Check if ticker is blacklisted
		for _, blacklistedTicker := range TickerBlacklist {
			if ticker == blacklistedTicker {
				blacklisted = true
			}
		}

		if blacklisted {
			continue
		}

		// Check for existing valid ticker
		for _, realTicker := range extractor.tickerList {
			if ticker == realTicker.Symbol {
				validTickers = append(validTickers, realTicker)
			}
		}
	}

	// If there are more than N tickers, count as spam and discard
	if len(validTickers) > TickerSpamThreshold {
		return []common.IEXTicker{}
	}

	return validTickers
}

// ExtractShortNames extracts a short name mentions from strings, given a list of tickers.
func (extractor *TickerExtractor) ExtractShortNames(s string, tickerList []common.IEXTicker) []common.IEXTicker {
	// Extract all alphanumeric words
	words := wordRegex.FindAllString(s, -1)

	// Find mentioned tickers
	mentionedTickers := []common.IEXTicker{}
	for _, ticker := range tickerList {

		if ticker.ShortName == "" {
			continue
		}

		for _, word := range words {
			if word == ticker.ShortName {
				mentionedTickers = append(mentionedTickers, ticker)
			}
		}
	}

	return mentionedTickers
}

// CalcTickerFrequency calculates frequency counts of each ticker and returns both
// a list of unique tickers and a map of type map[ticker.Symbol]frequencyCounts
func (extractor *TickerExtractor) CalcTickerFrequency(tickerList []common.IEXTicker) ([]common.IEXTicker, map[string]int) {
	frequencies := make(map[string]int)
	uniqueTickers := []common.IEXTicker{}

	for _, ticker := range tickerList {
		if frequencies[ticker.Symbol] == 0 {
			uniqueTickers = append(uniqueTickers, ticker)
		}
		frequencies[ticker.Symbol]++
	}

	return uniqueTickers, frequencies
}

// ExtractTickerMentions extracts every pb.TickerMention from a string
func (extractor *TickerExtractor) ExtractTickerMentions(s string, parentID string, parentSource string, timestamp *timestamppb.Timestamp) []pb.TickerMention {
	if s == "" {
		return []pb.TickerMention{}
	}

	// Extract tickers and shortname mentions from string
	tickers := extractor.ExtractTickers(s)
	nameTickers := extractor.ExtractShortNames(s, tickers)

	// Calculate frequencies
	uniqTickers, tickerFrequencies := extractor.CalcTickerFrequency(tickers)
	_, nameFrequencies := extractor.CalcTickerFrequency(nameTickers)

	// Calculate extra metrics
	questionCount := len(questionRegex.FindAllString(s, -1))
	wordCount := len(wordRegex.FindAllString(s, -1))

	// Generate list of ticker mentions
	mentions := []pb.TickerMention{}
	for _, ticker := range uniqTickers {
		mentions = append(mentions, pb.TickerMention{
			Symbol:            ticker.Symbol,
			ParentId:          parentID,
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
