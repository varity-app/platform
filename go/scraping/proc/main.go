package main

import (
	"log"

	"github.com/VarityPlatform/scraping/common"
	"google.golang.org/protobuf/types/known/timestamppb"

	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
)

// Entrypoint method
func main() {
	db := common.InitPostgres()

	allTickers, err := fetchTickers(db)
	if err != nil {
		log.Fatalln("Error fetching tickers:", err.Error())
	}

	submission := redditPB.RedditSubmission{
		SubmissionId: "testid",
		Title:        "Buy Apple ($AAPL) stocks you cucks.  Also GME?",
		Body:         "This is the boring part. Oh well.  GOOGL.  I like Google?.",
		Timestamp:    timestamppb.Now(),
	}

	mentions := procRedditSubmission(&submission, allTickers)

	for _, mention := range mentions {
		log.Println(mention.ParentId, mention.Symbol, mention.SymbolCounts, mention.ShortNameCounts, mention.QuestionMarkCount, mention.WordCount)
	}
}
