package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/VarityPlatform/scraping/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SubmissionsIntervalMinutes is the interval in minutes at which submissions will be scraped
const SubmissionsIntervalMinutes int = 5

// CommentsIntervalMinutes is the interval in minutes at which comments will be scraped
const CommentsIntervalMinutes int = 1

// NumWorkers is the number of parallel workers that will create new Cloud Tasks
const NumWorkers int = 10

// Limit is the maximum number of posts to scrape from one request.  100 is the max value allowed by PSAW.
const Limit int = 100

// Entrypoint method
func main() {
	// Initialize task submitter
	ctx := context.Background()

	// Initialize cobra commands
	submissionsCmd := &cobra.Command{
		Use:   "submissions [start date] [end date]",
		Short: "Scrape historical submissions from PSAW",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			subreddit := viper.GetString("subreddit")
			deployment := viper.GetString("deployment")
			url := viper.GetString("url")

			// Check for valid deployment flag
			if deployment != common.DeploymentModeDev && deployment != common.DeploymentModeProd {
				log.Fatalf("invalid deployment mode: '%s'.  Must be one of (dev, prod)", deployment)
			}

			// Init task submitter
			svcEmail := fmt.Sprintf("varity-tasks-svc-%s@varity.iam.gserviceaccount.com", deployment)
			queueID := fmt.Sprintf("reddit-historical-%s", deployment)
			ts, err := NewTaskSubmitter(ctx, svcEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer ts.Close()

			// Parse date arguments
			start, err := parseDate(args[0])
			if err != nil {
				log.Fatalf("invalid start date: %s.", args[0])
			}
			end, err := parseDate(args[1])
			if err != nil {
				log.Fatalf("invalid end date: %s.", args[1])
			}

			// Generate list of ranges
			befores, afters := generateDates(start, end, SubmissionsIntervalMinutes)

			// Create worker threads
			bodiesChan := make(chan *requestBody)
			wg := new(sync.WaitGroup)
			for i := 0; i < NumWorkers; i++ {
				go func() {
					for body := range bodiesChan {
						err = ts.Create(ctx, queueID, url+"/scraping/reddit/historical/submissions", body)
						if err != nil {
							log.Fatal(err)
						}
						wg.Done()
					}
				}()
			}

			// Create bodies
			for i := 0; i < len(befores); i++ {
				body := newRequestBody(subreddit, befores[i], afters[i], Limit)
				wg.Add(1)
				bodiesChan <- body
			}
			wg.Wait()

			fmt.Printf("Created %d tasks!", len(befores))

		},
	}

	commentsCmd := &cobra.Command{
		Use:   "comments [start date] [end date]",
		Short: "Scrape historical comments from PSAW",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			subreddit := viper.GetString("subreddit")
			deployment := viper.GetString("deployment")
			url := viper.GetString("url")

			// Check for valid deployment flag
			if deployment != common.DeploymentModeDev && deployment != common.DeploymentModeProd {
				log.Fatalf("Invalid deployment mode: '%s'.  Must be one of (dev, prod)", deployment)
			}

			// Init task submitter
			svcEmail := fmt.Sprintf("varity-tasks-svc-%s@varity.iam.gserviceaccount.com", deployment)
			queueID := fmt.Sprintf("reddit-historical-%s", deployment)
			ts, err := NewTaskSubmitter(ctx, svcEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer ts.Close()

			// Parse date arguments
			start, err := parseDate(args[0])
			if err != nil {
				log.Printf("Error: invalid start date: %s.", args[0])
			}
			end, err := parseDate(args[1])
			if err != nil {
				log.Printf("Error: invalid end date: %s.", args[1])
			}

			// Generate list of ranges
			befores, afters := generateDates(start, end, CommentsIntervalMinutes)

			// Create worker threads
			bodiesChan := make(chan *requestBody)
			wg := new(sync.WaitGroup)
			for i := 0; i < NumWorkers; i++ {
				go func() {
					for body := range bodiesChan {
						err = ts.Create(ctx, queueID, url+"/scraping/reddit/historical/comments", body)
						if err != nil {
							log.Fatal(err)
						}
						wg.Done()
					}
				}()
			}

			// Create bodies
			for i := 0; i < len(befores); i++ {
				body := newRequestBody(subreddit, befores[i], afters[i], Limit)
				wg.Add(1)
				bodiesChan <- body
			}
			wg.Wait()

			fmt.Printf("Created %d tasks!", len(befores))

		},
	}

	// Create rootCmd
	rootCmd := &cobra.Command{}

	// Create flags
	rootCmd.PersistentFlags().StringP("subreddit", "s", "", "subreddit to scrape")
	rootCmd.PersistentFlags().StringP("deployment", "d", "dev", "deployment to publish to")
	rootCmd.PersistentFlags().StringP("url", "u", "", "url to publish to")
	viper.BindPFlag("subreddit", rootCmd.PersistentFlags().Lookup("subreddit"))
	viper.BindPFlag("deployment", rootCmd.PersistentFlags().Lookup("deployment"))
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	if err := rootCmd.MarkPersistentFlagRequired("subreddit"); err != nil {
		log.Fatal(err)
	} else if err := rootCmd.MarkPersistentFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	// Run cobra commands
	rootCmd.AddCommand(submissionsCmd)
	rootCmd.AddCommand(commentsCmd)
	rootCmd.Execute()
}
