package bigquery2influx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/labstack/echo/v4"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data"
	b2i "github.com/varity-app/platform/scraping/internal/data/bigquery2influx"
	transforms "github.com/varity-app/platform/scraping/internal/transforms/mentions"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/logger"
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
	Influx         InfluxOpts
	Postgres       data.PostgresOpts
	DeploymentMode string
}

// NewService creates a new Bigquery->InfluxDB service in the form of an Echo server.
func NewService(ctx context.Context, opts ServiceOpts) (*echo.Echo, error) {

	// Init bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	// Init InfluxDB client
	influxClient := influxdb2.NewClient(opts.Influx.Addr, opts.Influx.Token)

	// Init postgres client
	escapedUsername := url.QueryEscape(opts.Postgres.Username)
	escapedPassword := url.QueryEscape(opts.Postgres.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", escapedUsername, escapedPassword, opts.Postgres.Address, opts.Postgres.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres.NewDriver: %v", err)
	}

	// Init repos
	membershipRepo := b2i.NewCohortMembershipRepo(bqClient, db, opts.DeploymentMode)
	mentionsRepo := b2i.NewMentionsRepo(bqClient, opts.DeploymentMode)

	// Init Echo
	web := echo.New()
	web.HideBanner = true
	registerRoutes(web, influxClient, membershipRepo, mentionsRepo, opts)

	return web, nil
}

// Register API routes
func registerRoutes(web *echo.Echo, influxClient influxdb2.Client, membershipRepo *b2i.CohortMembershipRepo, mentionsRepo *b2i.MentionsRepo, opts ServiceOpts) {

	web.POST("/api/v1/etl/bigquery-to-influx", func(c echo.Context) error {

		// Parse request body
		var request etlv1.PubSubRequest
		if err := c.Bind(&request); err != nil {
			log.Printf("echo.BindRequestBody: %v", err)
			return echo.ErrInternalServerError
		}

		// Decode base64 encoded request spec
		decoded, err := base64.StdEncoding.DecodeString(request.Message.Data)
		if err != nil {
			log.Printf("base64.DecodeString: %v", err)
			return echo.ErrInternalServerError
		}

		// Parse request spec decoded JSON string
		var spec etlv1.BigqueryToInfluxSpec
		if err := json.Unmarshal(decoded, &spec); err != nil {
			log.Println(request.Message.Data)
			log.Printf("json.UnmarshalSpec: %v", err)
			return echo.ErrBadRequest
		}

		// Get request context
		ctx := c.Request().Context()

		// Fetch cohort memberships from the past month
		memberships, err := membershipRepo.GetMemberships(ctx, spec.Year, spec.Month-1)
		if err != nil {
			log.Printf("membershipRepo.GetMemberships: %v", err)
			return echo.ErrInternalServerError
		}

		// Fetch mentions for the given hour
		mentions, err := mentionsRepo.GetMentions(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			log.Printf("mentionsRepo.GetMentions: %v", err)
			return echo.ErrInternalServerError
		}

		// Aggregate mentions
		ts := time.Date(spec.Year, time.Month(spec.Month), spec.Day, spec.Hour, 0, 0, 0, time.UTC)
		points := transforms.Aggregate(memberships, mentions, ts)

		// Write points to InfluxDB
		writeAPI := influxClient.WriteAPIBlocking(opts.Influx.Org, opts.Influx.Bucket)

		if err := writeAPI.WritePoint(ctx, points...); err != nil {
			log.Printf("influxdb.WritePoint: %v", err)
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, etlv1.BigqueryToInfluxResponse{PointsCount: len(points)})
	})

}
