package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/VarityPlatform/scraping/common"
)

type requestBody struct {
	Subreddit string `json:"subreddit"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Limit     int    `json:"limit"`
}

func newRequestBody(subreddit string, before time.Time, after time.Time, limit int) *requestBody {
	return &requestBody{
		Subreddit: subreddit,
		Before:    before.Format(time.RFC3339),
		After:     after.Format(time.RFC3339),
		Limit:     limit,
	}
}

// TaskSubmitter submits tasks to the Cloud Tasks API
type TaskSubmitter struct {
	client *cloudtasks.Client
	email  string // Service account email address
}

// NewTaskSubmitter creates a new TaskSubmitter instance
func NewTaskSubmitter(ctx context.Context, email string) (*TaskSubmitter, error) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("taskSubmitter.NewClient: %v", err)
	}

	return &TaskSubmitter{
		client: client,
		email:  email,
	}, nil
}

// Close closes a TaskSubmitter's connection
func (ts *TaskSubmitter) Close() error {
	return ts.client.Close()
}

// Create creates and submits a new task to the queue
func (ts *TaskSubmitter) Create(ctx context.Context, queueID string, url string, body *requestBody) error {

	// Serialize the request body.
	serializedBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("taskSubmitter.SerializeBody: %v", err)
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", common.GCPProjectID, common.GCPRegion, queueID)

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: ts.email,
						},
					},
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: serializedBody,
				},
			},
		},
	}

	_, err = ts.client.CreateTask(ctx, req)
	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return nil
}
