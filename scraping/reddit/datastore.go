package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/spf13/viper"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Take a list of submissions and check if their IDs exist in datastore.
// Return a list of submissions that do not currently exist
func getNewDSSubmissions(ctx context.Context, fsClient *firestore.Client, posts []*reddit.Post) ([]*reddit.Post, error) {

	// Create list of firestore refs for each submission
	refs := [](*firestore.DocumentRef){}
	collection := fsClient.Collection(REDDIT_SUBMISSIONS + "-" + viper.GetString("deploymentMode"))

	for _, post := range posts {
		refs = append(refs, collection.Doc(post.ID))
	}

	// Fetch submissions from firestore
	docsnaps, err := fsClient.GetAll(ctx, refs)
	if err != nil {
		return nil, fmt.Errorf("firestore.GetAll: %v", err)
	}

	// Create batch of writes
	newPosts := []*reddit.Post{}

	// Check each snapshot for each document
	for idx, snap := range docsnaps {
		// If there is no entry for a submission, then create one
		if !snap.Exists() {
			newPosts = append(newPosts, posts[idx])
		}
	}

	return newPosts, nil
}

// Save a list of comments to datastore
func saveNewDSSubmissions(ctx context.Context, fsClient *firestore.Client, submissions []*reddit.Post) error {
	collection := fsClient.Collection(REDDIT_SUBMISSIONS + "-" + viper.GetString("deploymentMode"))

	// Create batch of writes
	batch := fsClient.Batch()
	willCommit := false

	// Check each snapshot for each document
	for _, submission := range submissions {
		batch.Create(collection.Doc(submission.ID), map[string]interface{}{})
		willCommit = true
	}

	// There are changes to write
	if willCommit {
		_, err := batch.Commit(ctx)
		if err != nil {
			return fmt.Errorf("firestore.WriteBatch: %v", err)
		}
	}

	return nil
}

// Take a list of comments and check if their IDs exist in datastore.
// Return a list of comments that do not currently exist
func getNewDSComments(ctx context.Context, fsClient *firestore.Client, comments []*reddit.Comment) ([]*reddit.Comment, error) {

	// Create list of firestore refs for each comment
	refs := [](*firestore.DocumentRef){}
	collection := fsClient.Collection(REDDIT_COMMENTS + "-" + viper.GetString("deploymentMode"))

	for _, comment := range comments {
		refs = append(refs, collection.Doc(comment.ID))
	}

	// Fetch comments from firestore
	docsnaps, err := fsClient.GetAll(ctx, refs)
	if err != nil {
		return nil, fmt.Errorf("firestore.GetAll: %v", err)
	}

	// Create batch of writes
	newComments := []*reddit.Comment{}

	// Check each snapshot for each document
	for idx, snap := range docsnaps {
		// If there is no entry for a comment, then create one
		if !snap.Exists() {
			newComments = append(newComments, comments[idx])
		}
	}

	return newComments, nil
}

// Save a list of comments to datastore
func saveNewDSComments(ctx context.Context, fsClient *firestore.Client, comments []*reddit.Comment) error {
	collection := fsClient.Collection(REDDIT_COMMENTS + "-" + viper.GetString("deploymentMode"))

	// Create batch of writes
	batch := fsClient.Batch()
	willCommit := false

	// Check each snapshot for each document
	for _, comment := range comments {
		batch.Create(collection.Doc(comment.ID), map[string]interface{}{})
		willCommit = true
	}

	// There are changes to write
	if willCommit {
		_, err := batch.Commit(ctx)
		if err != nil {
			return fmt.Errorf("firestore.WriteBatch: %v", err)
		}
	}

	return nil
}
