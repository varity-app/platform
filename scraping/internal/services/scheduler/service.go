package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/labstack/echo/v4"

	"github.com/varity-app/platform/scraping/internal/common"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
)

// PubSub topic to publish BigQuery->InfluxDB ETL requests to.
const pubsubB2ITopic = "etl-bigquery-to-influx"

// ServiceOpts stores configuration parameters for a Bigquery->InfluxDB service
type ServiceOpts struct {
	DeploymentMode string
}

// NewService creates a new Bigquery->InfluxDB service in the form of an Echo server.
func NewService(ctx context.Context, opts ServiceOpts) (*echo.Echo, error) {

	// Init bigquery client
	psClient, err := pubsub.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Init Echo
	web := echo.New()
	web.HideBanner = true
	registerRoutes(web, psClient, opts)

	return web, nil
}

// Register API routes
func registerRoutes(web *echo.Echo, psClient *pubsub.Client, opts ServiceOpts) {

	web.POST("/api/scheduler/v1/etl/bigquery-to-influx/recent", func(c echo.Context) error {

		// Get request context
		ctx := c.Request().Context()

		// Create bigquery to influx message spec
		now := time.Now()
		spec := etlv1.BigqueryToInfluxSpec{
			Year:  now.Year(),
			Month: int(now.Month()),
			Day:   now.Day(),
			Hour:  now.Hour() - 1,
		}

		// Serialize spec
		serializedSpec, err := json.Marshal(spec)
		if err != nil {
			log.Printf("json.Marshal: %v", err)
			return err
		}

		// Publish message
		topicName := fmt.Sprintf("%s-%s", pubsubB2ITopic, opts.DeploymentMode)
		log.Printf("Sending message: %s to %s", serializedSpec, topicName)
		topic := psClient.Topic(topicName)

		msg := pubsub.Message{Data: serializedSpec}
		if _, err = topic.Publish(ctx, &msg).Get(ctx); err != nil {
			log.Printf("pubsub.Publish: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{})
	})

}
