package main

import (
	"fmt"

	"github.com/varity-app/platform/scraping/internal/common"

	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// createTaskRequest constructs a Cloud Tasks request
func createTaskRequest(
	deploymentMode string,
	queueID string,
	endpoint string,
	body []byte,
) *taskspb.CreateTaskRequest {

	// Build service account email
	email := fmt.Sprintf(common.CloudTasksSvcAccount, deploymentMode)

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
