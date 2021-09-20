package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	reddit "github.com/varity-app/platform/scraping/api/reddit/v1"
	"github.com/varity-app/platform/scraping/internal/common"
)

const (
	// numCloudTaskWorkers is the number of parallel workers that will create new Cloud Tasks
	numCloudTaskWorkers int = 20

	// postsLimit is the maximum number of posts to scrape from one request.
	// 100 is the max value allowed by PSAW.
	postsLimit int = 100

	// submissionsIntervalMinutes is the interval in minutes at which submissions will be scraped
	submissionsIntervalMinutes int = 3

	// commentsIntervalMinutes is the interval in minutes at which comments will be scraped
	commentsIntervalMinutes int = 1

	// subredditWSB is a literal for "wallstreetbets"
	subredditWSB string = "wallstreetbets"
)

// Base reddit historical command
var historicalCmd = &cobra.Command{
	Use:   "historical",
	Short: "Scrape historical reddit posts from PSAW.",
}

// Scrape comments command
var historicalCommentsCmd = &cobra.Command{
	Use:   "comments [start date] [end date]",
	Short: "Scrape historical comments from PSAW",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Init cloud tasks client
		ctx := context.Background()
		ctClient, err := cloudtasks.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("pubsub.NewClient: %v", err)
		}
		defer ctClient.Close()

		// Parse flags
		subreddit, deployment, url, err := parseHistoricalFlags()
		if err != nil {
			return err
		}

		// Parse date arguments
		start, err := parseDate(args[0])
		if err != nil {
			return fmt.Errorf("invalid start date: %s", args[0])
		}
		end, err := parseDate(args[1])
		if err != nil {
			return fmt.Errorf("invalid end date: %s", args[1])
		}

		// Define cloud tasks queue
		queueID := fmt.Sprintf("%s-%s", common.CloudTasksQueueRedditHistorical, deployment)
		endpoint := url + "/scraping/reddit/historical/comments"

		// Generate list of ranges
		intervalMultiple := commentsIntervalMinutes
		if strings.ToLower(subreddit) != subredditWSB {
			intervalMultiple *= 5
		}
		if end.Year() < 2021 {
			intervalMultiple *= 2
		}
		if end.Year() < 2020 {
			intervalMultiple *= 2
		}
		starts, ends := generateMinuteRanges(start, end, intervalMultiple)

		// Create worker threads
		var errcList []<-chan error
		chunkSize := len(ends) / numCloudTaskWorkers
		for i := 0; i < len(starts); i += chunkSize {

			// Define chunks
			endsChunk := ends[i:min(i+chunkSize, len(ends))]
			startsChunk := starts[i:min(i+chunkSize, len(starts))]

			// Spawn thread
			errc := make(chan error, 1)
			go func(endsChunk, startsChunk []time.Time) {
				defer close(errc)

				for i := range startsChunk {
					start := startsChunk[i]
					end := endsChunk[i]

					// Create request body
					body := reddit.HistoricalRequest{
						Subreddit: subreddit,
						After:     start.Format(time.RFC3339),
						Before:    end.Format(time.RFC3339),
						Limit:     postsLimit,
					}

					// Serialize request body
					serializedBody, err := json.Marshal(body)
					if err != nil {
						errc <- fmt.Errorf("json.Marshal: %v", err)
						return
					}

					// Submit cloud task
					req := createTaskRequest(deployment, queueID, endpoint, serializedBody) // Submit cloud tasks request
					_, err = ctClient.CreateTask(ctx, req)
					if err != nil {
						errc <- fmt.Errorf("cloudtasks.CreateTask: %v", err)
						return
					}
				}
			}(endsChunk, startsChunk)
			errcList = append(errcList, errc)
		}

		// Wait for each thread to finish
		if err := common.WaitForPipeline(errcList...); err != nil {
			return err
		}

		fmt.Printf("Created %d historical reddit comment scraping tasks!", len(starts))
		return nil
	},
}

// Scrape submissions command
var historicalSubmissionsCmd = &cobra.Command{
	Use:   "submissions [start date] [end date]",
	Short: "Scrape historical submissions from PSAW",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Init cloud tasks client
		ctx := context.Background()
		ctClient, err := cloudtasks.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("pubsub.NewClient: %v", err)
		}
		defer ctClient.Close()

		// Parse flags
		subreddit, deployment, url, err := parseHistoricalFlags()
		if err != nil {
			return err
		}

		// Parse date arguments
		start, err := parseDate(args[0])
		if err != nil {
			return fmt.Errorf("invalid start date: %s", args[0])
		}
		end, err := parseDate(args[1])
		if err != nil {
			return fmt.Errorf("invalid end date: %s", args[1])
		}

		// Define cloud tasks queue
		queueID := fmt.Sprintf("%s-%s", common.CloudTasksQueueRedditHistorical, deployment)
		endpoint := url + "/scraping/reddit/historical/submissions"

		// Generate list of ranges
		intervalMultiple := submissionsIntervalMinutes
		if strings.ToLower(subreddit) != subredditWSB {
			intervalMultiple *= 5
		}
		if end.Year() < 2021 {
			intervalMultiple *= 2
		}
		if end.Year() < 2020 {
			intervalMultiple *= 2
		}
		starts, ends := generateMinuteRanges(start, end, intervalMultiple)

		// Create worker threads
		var errcList []<-chan error
		chunkSize := len(ends)/numCloudTaskWorkers + 1
		for i := 0; i < len(starts); i += chunkSize {

			// Define chunks
			endsChunk := ends[i:min(i+chunkSize, len(ends))]
			startsChunk := starts[i:min(i+chunkSize, len(starts))]

			// Spawn thread
			errc := make(chan error, 1)
			go func(endsChunk, startsChunk []time.Time) {
				defer close(errc)

				for i := range startsChunk {
					start := startsChunk[i]
					end := endsChunk[i]

					// Create request body
					body := reddit.HistoricalRequest{
						Subreddit: subreddit,
						After:     start.Format(time.RFC3339),
						Before:    end.Format(time.RFC3339),
						Limit:     postsLimit,
					}

					// Serialize request body
					serializedBody, err := json.Marshal(body)
					if err != nil {
						errc <- fmt.Errorf("json.Marshal: %v", err)
						return
					}

					// Submit cloud task
					req := createTaskRequest(deployment, queueID, endpoint, serializedBody) // Submit cloud tasks request
					_, err = ctClient.CreateTask(ctx, req)
					if err != nil {
						errc <- fmt.Errorf("cloudtasks.CreateTask: %v", err)
						return
					}
				}
			}(endsChunk, startsChunk)
			errcList = append(errcList, errc)
		}

		// Wait for each thread to finish
		if err := common.WaitForPipeline(errcList...); err != nil {
			return err
		}

		fmt.Printf("Created %d historical reddit submission scraping tasks!", len(starts))
		return nil
	},
}

// Parse flags for historical reddit commands
func parseHistoricalFlags() (string, string, string, error) {
	subreddit := strings.ToLower(viper.GetString("subreddit"))
	deployment := viper.GetString("deployment")
	url := viper.GetString("urls.reddit_historical")

	// Check for valid deployment flag
	if deployment != common.DeploymentModeDev && deployment != common.DeploymentModeProd {
		return "", "", "", fmt.Errorf(
			"invalid deployment mode: '%s'.  Must be one of (dev, prod)",
			deployment,
		)
	}

	return subreddit, deployment, url, nil
}

// Initialize the historical reddit scraping commands and add to the root command.
func initRedditHistorical(rootCmd *cobra.Command) {

	// Add subcommands to historical command
	historicalCmd.AddCommand(historicalCommentsCmd)
	historicalCmd.AddCommand(historicalSubmissionsCmd)

	// Add historical command to root
	rootCmd.AddCommand(historicalCmd)

	// Create flags
	historicalCmd.PersistentFlags().StringP("subreddit", "s", "", "subreddit to scrape")
	historicalCmd.PersistentFlags().StringP("url", "u", "", "url to publish to")
	viper.BindPFlag("subreddit", historicalCmd.PersistentFlags().Lookup("subreddit"))
	viper.BindPFlag("urls.reddit_historical", historicalCmd.PersistentFlags().Lookup("url"))

	if err := historicalCmd.MarkPersistentFlagRequired("subreddit"); err != nil {
		log.Fatal(err)
	} else if err := historicalCmd.MarkPersistentFlagRequired("url"); err != nil {
		log.Fatal(err)
	}
}

// Helper method for returning a minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
