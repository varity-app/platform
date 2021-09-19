package bigquery2influx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/labstack/echo/v4"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	bq "github.com/varity-app/platform/scraping/internal/data/bigquery"
	"github.com/varity-app/platform/scraping/internal/logging"
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
	Influx         InfluxOpts
	DeploymentMode string
}

// Service is microservice served via HTTP
type Service struct {
	bqClient *bigquery.Client
	web      *echo.Echo
	logger   *logging.Logger
}

// NewService creates a new Bigquery->InfluxDB service in the form of an Echo server.
func NewService(ctx context.Context, logger *logging.Logger, opts ServiceOpts) (*Service, error) {

	// Init bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	// Init InfluxDB client
	influxClient := influxdb2.NewClient(opts.Influx.Addr, opts.Influx.Token)

	// Init repos
	membershipRepo, err := bq.NewCohortMembershipRepo(bqClient, logger, opts.DeploymentMode)
	if err != nil {
		return nil, fmt.Errorf("cohortMembershipsRepo.New: %v", err)
	}

	submissionsRepo := bq.NewRedditSubmissionsRepo(bqClient, opts.DeploymentMode)
	commentsRepo := bq.NewRedditCommentsRepo(bqClient, opts.DeploymentMode)
	mentionsRepo := bq.NewMentionsRepo(bqClient, opts.DeploymentMode)

	// Init Echo
	web := echo.New()
	web.HideBanner = true

	service := &Service{
		bqClient: bqClient,
		web:      web,
		logger:   logger,
	}
	service.registerRoutes(web, influxClient, membershipRepo, submissionsRepo, commentsRepo, mentionsRepo, opts)

	return service, nil
}

// Close connections
func (s *Service) Close() error {
	return s.bqClient.Close()
}

// Start the service's HTTP web server
func (s *Service) Start(address string) {
	s.web.Logger.Fatal(s.web.Start(address))
}

// Register API routes
func (s *Service) registerRoutes(
	web *echo.Echo,
	influxClient influxdb2.Client,
	membershipRepo *bq.CohortMembershipRepo,
	submissionsRepo *bq.RedditSubmissionsRepo,
	commentsRepo *bq.RedditCommentsRepo,
	mentionsRepo *bq.MentionsRepo,
	opts ServiceOpts,
) {
	s.web.POST("/api/etl/bigquery-to-influx/v1/", func(c echo.Context) error {

		// Parse request body
		var spec etlv1.BigqueryToInfluxRequest
		if err := c.Bind(&spec); err != nil {
			log.Printf("echo.BindRequestBody: %v", err)
			return echo.ErrInternalServerError
		}

		// Get request context
		ctx := c.Request().Context()

		// Fetch mentions for the given hour
		mentions, err := mentionsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			s.logger.Error(fmt.Errorf("mentionsRepo.GetMentions: %v", err))
			return echo.ErrInternalServerError
		}
		s.logger.Debug(fmt.Sprintf("Fetched %d ticker mentions.", len(mentions)))

		// Fetch submissions for the given hour
		submissions, err := submissionsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			s.logger.Error(fmt.Errorf("submissionsRepo.GetMentions: %v", err))
			return echo.ErrInternalServerError
		}
		s.logger.Debug(fmt.Sprintf("Fetched %d submissions.", len(submissions)))

		// Fetch comments for the given hour
		comments, err := commentsRepo.Get(ctx, spec.Year, spec.Month, spec.Day, spec.Hour)
		if err != nil {
			s.logger.Error(fmt.Errorf("commentsRepo.GetMentions: %v", err))
			return echo.ErrInternalServerError
		}
		s.logger.Debug(fmt.Sprintf("Fetched %d comments.", len(comments)))

		// Augment tickers
		augmented := transforms.AugmentMentions(mentions, submissions, comments)
		s.logger.Debug(fmt.Sprintf("Created %d augmented ticker mentions.", len(augmented)))

		// Fetch cohort memberships from the past month
		memberships, err := membershipRepo.Get(ctx, spec.Year, spec.Month-1)
		if err != nil {
			s.logger.Error(fmt.Errorf("membershipRepo.GetMemberships: %v", err))
			return echo.ErrInternalServerError
		}
		s.logger.Debug(fmt.Sprintf("Fetched %d cohort memberships.", len(memberships)))

		// Aggregate mentions
		ts := time.Date(spec.Year, time.Month(spec.Month), spec.Day, spec.Hour, 0, 0, 0, time.UTC)
		points := transforms.Aggregate(memberships, augmented, ts)

		// Write points to InfluxDB
		writeAPI := influxClient.WriteAPIBlocking(opts.Influx.Org, opts.Influx.Bucket)

		if err := writeAPI.WritePoint(ctx, points...); err != nil {
			s.logger.Error(fmt.Errorf("influxdb.WritePoint: %v", err))
			return echo.ErrInternalServerError
		}
		s.logger.Debug(fmt.Sprintf("Wrote %d points to InfluxDB.", len(points)))

		return c.JSON(http.StatusOK, etlv1.BigqueryToInfluxResponse{PointsCount: len(points)})
	})

}
