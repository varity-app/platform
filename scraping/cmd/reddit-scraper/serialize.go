package main

import (
	"fmt"

	rpb "github.com/varity-app/platform/scraping/api/reddit/v1"
	"google.golang.org/protobuf/proto"
)

// Serialize a list of reddit submissions into a list of byte arrays
func serializeSubmissions(submissions []*rpb.RedditSubmission) ([][]byte, error) {
	serializedSubmissions := [][]byte{}
	for _, submission := range submissions {
		serializedSubmission, err := proto.Marshal(submission)
		if err != nil {
			return nil, fmt.Errorf("serialize.Submissions: %v", err)
		}
		serializedSubmissions = append(serializedSubmissions, serializedSubmission)
	}

	return serializedSubmissions, nil
}

// Serialize a list of reddit comments into a list of byte arrays
func serializeComments(comments []*rpb.RedditComment) ([][]byte, error) {
	serializedComments := [][]byte{}
	for _, comment := range comments {
		serializedComment, err := proto.Marshal(comment)
		if err != nil {
			return nil, fmt.Errorf("serialize.Comments: %v", err)
		}
		serializedComments = append(serializedComments, serializedComment)
	}

	return serializedComments, nil
}
