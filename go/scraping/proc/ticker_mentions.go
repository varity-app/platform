package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"

	"github.com/VarityPlatform/scraping/common"
	commonPB "github.com/VarityPlatform/scraping/protobuf/common"
	"google.golang.org/protobuf/proto"

	"github.com/spf13/viper"

	"go.opentelemetry.io/otel/trace"
)

// Read ticker mentions from Pub/Sub and process them.
func readTickerMentions(ctx context.Context, psClient *pubsub.Client, bqClient *bigquery.Client, tracer *trace.Tracer) (int, error) {

	sub := psClient.Subscription(common.SUBSCRIPTION_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))
	topic := psClient.Topic(common.TOPIC_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))

	cm := make(chan *pubsub.Message)
	defer close(cm)

	readWG := new(sync.WaitGroup)
	writeWG := new(sync.WaitGroup)
	var mu sync.Mutex

	// Count messages processed
	count := 0
	totalCount := 0

	// Create bigquery streaming insert
	inserter := bqClient.Dataset(
		common.BIGQUERY_DATASET_SCRAPING + "_" + viper.GetString("deploymentMode"),
	).Table(common.BIGQUERY_TABLE_TICKER_MENTIONS).Inserter()

	// Async errors
	var readErr error = nil
	var writeErr error = nil

	// Spawn thread for parsing protobuf messages
	mentions := []*commonPB.TickerMention{}
	go func() {
		for msg := range cm {
			mention := &commonPB.TickerMention{}
			if err := proto.Unmarshal(msg.Data, mention); err != nil {
				log.Println(err)
				readErr = fmt.Errorf("tickerMention.Unmarshal: %v", err)
				readWG.Done()
				return
			}
			mu.Lock()
			mentions = append(mentions, mention)
			mu.Unlock()

			msg.Ack()
			readWG.Done()
		}
	}()

	// Receive messages for N sec
	for {
		pubCtx, cancel := context.WithTimeout(ctx, WAIT_INTERVAL)
		defer cancel()

		// Reset count
		count = 0
		mentions = []*commonPB.TickerMention{}

		ctx, span := (*tracer).Start(ctx, "pubsub.Receive")
		// Receive blocks until the context is cancelled or an error occurs.
		err := sub.Receive(pubCtx, func(ctx context.Context, msg *pubsub.Message) {
			if readErr != nil {
				cancel()
				return
			}
			readWG.Add(1)
			cm <- msg
			count++
		})
		if err != nil {
			log.Println(err)
			return totalCount, err
		}

		readWG.Wait()
		if readErr != nil {
			return totalCount, readErr
		}
		span.End()

		ctx, span = (*tracer).Start(ctx, "bigquery.Upload")
		// Upload to bigquery
		if err := inserter.Put(ctx, mentions); err != nil {
			// Republish mentions to pubsub if there's an error
			for _, mention := range mentions {

				serializedMention, serErr := proto.Marshal(mention)
				if serErr != nil {
					return totalCount, fmt.Errorf("mention.SerializeProto: %v AND bigquery.Insert: %v", serErr, err)
				}

				// Asynchronously publish to pubsub
				result := topic.Publish(ctx, &pubsub.Message{
					Data: serializedMention,
				})

				writeWG.Add(1)
				go func(res *pubsub.PublishResult) {
					defer writeWG.Done()
					_, err := result.Get(ctx)
					if err != nil {
						mu.Lock()
						writeErr = fmt.Errorf("mention.Publish: %v", err)
						mu.Unlock()
					}
					writeWG.Done()
				}(result)
			}

			writeWG.Wait()

			// Handle
			if writeErr != nil {
				return totalCount, fmt.Errorf("pubsub.PublishResult: %v AND bigquery.Insert: %v", writeErr, err)
			}
			return totalCount, fmt.Errorf("bigquery.Insert: %v", err)
		}
		span.End()

		totalCount += count

		// If no messages were received in this window, break
		if count == 0 {
			break
		}
	}

	return totalCount, nil

}
