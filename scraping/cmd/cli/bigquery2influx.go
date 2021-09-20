package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
	"github.com/varity-app/platform/scraping/internal/common"
)

var b2iCmd = &cobra.Command{
	Use:   "b2i [start date] [end date]",
	Short: "Create Cloud Tasks for exporting ticker mentions from BigQuery to InfluxDB",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Parse flags
		url := viper.GetString("urls.b2i")
		deployment := viper.GetString("deployment")

		// Check for valid deployment flag
		if deployment != common.DeploymentModeDev && deployment != common.DeploymentModeProd {
			return fmt.Errorf(
				"invalid deployment mode: '%s'.  Must be one of (dev, prod)",
				deployment,
			)
		}

		log.Println(url)

		// Init cloud tasks client
		ctx := context.Background()
		ctClient, err := cloudtasks.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("pubsub.NewClient: %v", err)
		}
		defer ctClient.Close()

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
		queueID := fmt.Sprintf("%s-%s", common.CloudTasksQueueB2I, deployment)
		endpoint := url + "/api/etl/bigquery-to-influx/v1/"

		// Generate list of ranges
		starts, _ := generateHourRanges(start, end)

		// Create worker threads
		var errcList []<-chan error
		chunkSize := len(starts) / numCloudTaskWorkers
		for i := 0; i < len(starts); i += chunkSize {

			// Define chunks
			startsChunk := starts[i:min(i+chunkSize, len(starts))]

			// Spawn thread
			errc := make(chan error, 1)
			go func(startsChunk []time.Time) {
				defer close(errc)

				for i := range startsChunk {
					start := startsChunk[i]

					// Create bigquery to influx message spec
					spec := etlv1.BigqueryToInfluxRequest{
						Year:  start.Year(),
						Month: int(start.Month()),
						Day:   start.Day(),
						Hour:  start.Hour(),
					}

					// Serialize spec
					serializedSpec, err := json.Marshal(spec)
					if err != nil {
						errc <- fmt.Errorf("json.Marshal: %v", err)
						return
					}

					// Submit cloud task
					req := createTaskRequest(deployment, queueID, endpoint, serializedSpec) // Submit cloud tasks request
					_, err = ctClient.CreateTask(ctx, req)
					if err != nil {
						errc <- fmt.Errorf("cloudtasks.CreateTask: %v", err)
						return
					}
				}
			}(startsChunk)
			errcList = append(errcList, errc)
		}

		// Wait for each thread to finish
		if err := common.WaitForPipeline(errcList...); err != nil {
			return err
		}

		fmt.Printf("Created %d BigQuery -> InfluxDB ETL tasks!", len(starts))
		return nil
	},
}

// Initialize the bigquery2influx command and add to the root command.
func initB2I(rootCmd *cobra.Command) {

	// Add bigquery2influx command to root command
	rootCmd.AddCommand(b2iCmd)

	// Set flags
	b2iCmd.Flags().StringP("url", "u", "", "url to publish to")
	if err := viper.BindPFlag("urls.b2i", b2iCmd.Flags().Lookup("url")); err != nil {
		log.Fatal(err)
	}

	if err := b2iCmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}
}
