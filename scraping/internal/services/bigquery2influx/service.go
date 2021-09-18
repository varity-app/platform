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
	bq "github.com/varity-app/platform/scraping/internal/data/bigquery"
	transforms "github.com/varity-app/platform/scraping/internal/transforms/mentions"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	membershipRepo := bq.NewCohortMembershipRepo(bqClient, db, opts.DeploymentMode)
	submissionsRepo := bq.NewRedditSubmissionsRepo(bqClient, opts.DeploymentMode)
	commentsRepo := bq.NewRedditCommentsRepo(bqClient, opts.DeploymentMode)
	mentionsRepo := bq.NewMentionsRepo(bqClient, opts.DeploymentMode)

	// Init Echo
	web := echo.New()
	web.HideBanner = true
	registerRoutes(web, influxClient, membershipRepo, submissionsRepo, commentsRepo, mentionsRepo, opts)

	return web, nil
}

// Register API routes
func registerRoutes(
	web *echo.Echo,
	influxClient influxdb2.Client,
	membershipRepo *bq.CohortMembershipRepo,
	submissionsRepo *bq.RedditSubmissionsRepo,
	commentsRepo *bq.RedditCommentsRepo,
	mentionsRepo *bq.MentionsRepo,
	opts ServiceOpts,
) {

	web.POST("/api/etl/bigquery-to-influx/v1/", func(c echo.Context) error {

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
		memberships, err := membershipRepo.Get(ctx, spec.Year, spec.Month-1)
		if err != nil {
			log.Printf("membershipRepo.GetMemberships: %v", err)
			return echo.ErrInternalServerError
		}

		// Fetch mentions for the given hour
		mentions, err := mentionsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			log.Printf("mentionsRepo.GetMentions: %v", err)
			return echo.ErrInternalServerError
		}

		// Fetch submissions for the given hour
		submissions, err := submissionsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			log.Printf("submissionsRepo.GetMentions: %v", err)
			return echo.ErrInternalServerError
		}

		// Fetch comments for the given hour
		comments, err := commentsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			log.Printf("commentsRepo.GetMentions: %v", err)
			return echo.ErrInternalServerError
		}

		// Augment tickers
		augmented := transforms.AugmentMentions(mentions, submissions, comments)

		// Aggregate mentions
		ts := time.Date(spec.Year, time.Month(spec.Month), spec.Day, spec.Hour, 0, 0, 0, time.UTC)
		points := transforms.Aggregate(memberships, augmented, ts)

		// Write points to InfluxDB
		writeAPI := influxClient.WriteAPIBlocking(opts.Influx.Org, opts.Influx.Bucket)

		if err := writeAPI.WritePoint(ctx, points...); err != nil {
			log.Printf("influxdb.WritePoint: %v", err)
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, etlv1.BigqueryToInfluxResponse{PointsCount: len(points)})
	})

}
