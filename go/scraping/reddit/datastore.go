package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"

	"github.com/spf13/viper"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Take a list of submissions and check if their IDs exist in datastore.
// Return a list of submissions that do not currently exist
func getNewDSSubmissions(ctx context.Context, fsClient *firestore.Client, posts []*reddit.Post) []*reddit.Post {

	// Create list of firestore refs for each submission
	refs := [](*firestore.DocumentRef){}
	collection := fsClient.Collection(REDDIT_SUBMISSIONS + "-" + viper.GetString("deploymentMode"))

	for _, post := range posts {
		refs = append(refs, collection.Doc(post.ID))
	}

	// Fetch submissions from firestore
	docsnaps, err := fsClient.GetAll(ctx, refs)
	if err != nil {
		log.Fatal("Error fetching from firestore:", err.Error())
	}

	// Create batch of writes
	batch := fsClient.Batch()
	willCommit := false
	newPosts := []*reddit.Post{}

	// Check each snapshot for each document
	for idx, snap := range docsnaps {
		// If there is no entry for a submission, then create one
		if !snap.Exists() {
			batch.Create(snap.Ref, map[string]interface{}{})
			willCommit = true
			newPosts = append(newPosts, posts[idx])
		}
	}

	// There are changes to write
	if willCommit {
		_, err = batch.Commit(ctx)
		if err != nil {
			log.Fatal("Error commiting firestore batch", err)
		}
	}

	return newPosts
}

// Take a list of comments and check if their IDs exist in datastore.
// Return a list of comments that do not currently exist
func getNewDSComments(ctx context.Context, fsClient *firestore.Client, comments []*reddit.Comment) []*reddit.Comment {

	// Create list of firestore refs for each comment
	refs := [](*firestore.DocumentRef){}
	collection := fsClient.Collection(REDDIT_COMMENTS + "-" + viper.GetString("deploymentMode"))

	for _, comment := range comments {
		refs = append(refs, collection.Doc(comment.ID))
	}

	// Fetch comments from firestore
	docsnaps, err := fsClient.GetAll(ctx, refs)
	if err != nil {
		log.Fatal("Error fetching from firestore:", err.Error())
	}

	// Create batch of writes
	batch := fsClient.Batch()
	willCommit := false
	newComments := []*reddit.Comment{}

	// Check each snapshot for each document
	for idx, snap := range docsnaps {
		// If there is no entry for a comment, then create one
		if !snap.Exists() {
			batch.Create(snap.Ref, map[string]interface{}{})
			willCommit = true
			newComments = append(newComments, comments[idx])
		}
	}

	// There are changes to write
	if willCommit {
		_, err = batch.Commit(ctx)
		if err != nil {
			log.Fatal("Error commiting firestore batch", err)
		}
	}

	return newComments
}
