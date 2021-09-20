package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/labstack/echo/v4"

	etlv1 "github.com/varity-app/platform/scraping/api/etl/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/logging"
)

var logger = logging.NewLogger()

// ServiceOpts stores configuration parameters for a Bigquery->InfluxDB service
type ServiceOpts struct {
	DeploymentMode string
	B2IUrl         string
}

// Service is microservice served via HTTP
type Service struct {
	ctClient       *cloudtasks.Client
	web            *echo.Echo
	deploymentMode string
}

// NewService creates a new Bigquery->InfluxDB service in the form of an Echo server.
func NewService(ctx context.Context, opts ServiceOpts) (*Service, error) {

	// Init cloud tasks client
	ctClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Init Echo
	web := echo.New()
	web.HideBanner = true

	// Create service instance
	service := &Service{
		ctClient:       ctClient,
		web:            web,
		deploymentMode: opts.DeploymentMode,
	}

	// Register routes
	service.registerRoutes(opts.B2IUrl)

	return service, nil
}

// Close connections
func (s *Service) Close() error {
	return s.ctClient.Close()
}

// Start the service's HTTP web server
func (s *Service) Start(address string) {
	s.web.Logger.Fatal(s.web.Start(address))
}

// Register API routes
func (s *Service) registerRoutes(b2iURL string) {

	s.web.POST("/api/scheduler/v1/etl/bigquery-to-influx/recent", func(c echo.Context) error {

		// Get request context
		ctx := c.Request().Context()

		// Create bigquery to influx message spec
		ts := time.Now().Add(-time.Hour)
		spec := etlv1.BigqueryToInfluxRequest{
			Year:  ts.Year(),
			Month: int(ts.Month()),
			Day:   ts.Day(),
			Hour:  ts.Hour(),
		}

		// Serialize spec
		serializedSpec, err := json.Marshal(spec)
		if err != nil {
			logger.Error(fmt.Errorf("json.Marshal: %v", err))
			return err
		}

		// Create cloud tasks request
		queueID := fmt.Sprintf("%s-%s", common.CloudTasksQueueB2I, s.deploymentMode)
		b2iEndpoint := b2iURL + "/api/etl/bigquery-to-influx/v1/"
		req := s.createTaskRequest(queueID, b2iEndpoint, serializedSpec)
		logger.Debug(fmt.Sprintf("Submitting task that will send to `%s`.", b2iEndpoint))

		// Submit cloud tasks request
		_, err = s.ctClient.CreateTask(ctx, req)
		if err != nil {
			logger.Error(fmt.Errorf("cloudtasks.CreateTask: %v", err))
			return echo.ErrInternalServerError
		}
		logger.Debug(fmt.Sprintf("Successfully submitted task with spec: %s", serializedSpec))

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Task created successfully!",
		})
	})

}

// createTaskRequest constructs a full Cloud Tasks request
func (s *Service) createTaskRequest(queueID string, endpoint string, body []byte) *taskspb.CreateTaskRequest {

	// Build service account email
	email := fmt.Sprintf(common.CloudTasksSvcAccount, s.deploymentMode)

	// Build the Task queue path.
	queuePath := fmt.Sprintf(
		"projects/%s/locations/%s/queues/%s",
		common.GCPProjectID,
		common.GCPRegion,
		queueID,
	)

	// Build the request
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        endpoint,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: email,
						},
					},
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: body,
				},
			},
		},
	}

	return req

}
