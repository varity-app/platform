package bigquery2influx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/go-redis/redis/v8"

	"github.com/labstack/echo/v4"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	b2i "github.com/varity-app/platform/scraping/internal/data/bigquery2influx"
	transforms "github.com/varity-app/platform/scraping/internal/transforms/mentions"
)

// InfluxOpts stores configuration parameters for an InfluxDB Client
type InfluxOpts struct {
	Addr   string
	Token  string
	Org    string
	Bucket string
}

// ServiceOpts stores configuration parameters for a Bigquery->InfluxDB service
type ServiceOpts struct {
	Redis          *redis.Options
	Influx         InfluxOpts
	DeploymentMode string
}

// NewService creates a new Bigquery->InfluxDB service in the form of an Echo server.
func NewService(ctx context.Context, opts ServiceOpts) (*echo.Echo, error) {

	// Init bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	// Init InfluxDB client
	influxClient := influxdb2.NewClient(opts.Influx.Addr, opts.Influx.Token)
	defer influxClient.Close()

	// Init repos
	membershipRepo := b2i.NewCohortMembershipRepo(bqClient, opts.Redis)
	mentionsRepo := b2i.NewMentionsRepo(bqClient, opts.Redis)

	// Init Echo
	web := echo.New()
	web.HideBanner = true
	registerRoutes(web, influxClient, membershipRepo, mentionsRepo, opts)

	return web, nil
}

// Register API routes
func registerRoutes(web *echo.Echo, influxClient influxdb2.Client, membershipRepo *b2i.CohortMembershipRepo, mentionsRepo *b2i.MentionsRepo, opts ServiceOpts) {

	web.POST("/api/v1/bigquery-to-influx", func(c echo.Context) error {

		// Parse request body
		var request etlv1.PubSubRequest
		if err := c.Bind(&request); err != nil {
			log.Printf("echo.BindRequestBody: %v", err)
			return err
		}

		// Parse request spec from Pub/Sub message body
		var spec etlv1.BigqueryToInfluxSpec
		if err := json.Unmarshal([]byte(request.Message.Data), &spec); err != nil {
			log.Printf("json.UnmarshalSpec: %v", err)
			return echo.ErrBadRequest
		}

		// Get request context
		ctx := c.Request().Context()

		// Fetch cohort memberships from the past month
		tableName := wrapDBTTable(opts.DeploymentMode, b2i.CohortTableRedditorsAge)
		memberships, err := membershipRepo.GetMemberships(ctx, tableName, spec.Year, spec.Month-1)
		if err != nil {
			log.Printf("membershipRepo.GetMemberships: %v", err)
			return err
		}

		// Fetch mentions for the given hour
		tableName = wrapDBTTable(opts.DeploymentMode, b2i.MentionsTableReddit)
		mentions, err := mentionsRepo.GetMentions(ctx, tableName, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			log.Printf("mentionsRepo.GetMentions: %v", err)
			return err
		}

		// Aggregate mentions
		ts := time.Date(spec.Year, time.Month(spec.Month), spec.Day, spec.Hour, 0, 0, 0, time.UTC)
		points := transforms.Aggregate(memberships, mentions, ts)

		// Write points to InfluxDB
		writeAPI := influxClient.WriteAPIBlocking(opts.Influx.Org, opts.Influx.Bucket)

		if err := writeAPI.WritePoint(ctx, points...); err != nil {
			log.Printf("influxdb.WritePoint: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, etlv1.BigqueryToInfluxResponse{PointsCount: len(points)})
	})

}

// Wrap a DBT Bigquery table name with the GCP project and dataset name.
func wrapDBTTable(deployment, tableName string) string {
	return fmt.Sprintf(
		"%s.%s_%s.%s",
		common.GCPProjectID,
		common.BigqueryDatasetDBTScraping,
		deployment,
		tableName,
	)
}
