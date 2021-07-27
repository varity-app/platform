package main

import (
	pb "github.com/VarityPlatform/scraping/protobuf/reddit"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vartanbeno/go-reddit/reddit"
)

// Encode reddit submission as protobuf
func submissionToProto(post *reddit.Post) *pb.RedditSubmission {

	submissionProto := &pb.RedditSubmission{
		SubmissionId: post.ID,
		Subreddit:    post.SubredditName,
		Title:        &post.Title,
		CreatedAt:    timestamppb.New(post.Created.Time),
		Body:         &post.Body,
		Author:       post.Author,
		AuthorId:     post.AuthorID,
		IsSelf:       &post.IsSelfPost,
		Permalink:    &post.Permalink,
		Url:          &post.URL,
	}

	return submissionProto
}

// Encode reddit comment as protobuf
func commentToProto(comment *reddit.Comment) *pb.RedditComment {

	commentProto := &pb.RedditComment{
		CommentId:    comment.ID,
		SubmissionId: comment.ParentID,
		Subreddit:    comment.SubredditName,
		CreatedAt:    timestamppb.New(comment.Created.Time),
		Body:         &comment.Body,
		Author:       comment.Author,
		AuthorId:     comment.AuthorID,
		Permalink:    &comment.Permalink,
	}

	return commentProto
}
