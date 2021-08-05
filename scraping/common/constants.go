package common

// GCPProjectID is the google cloud project name
const GCPProjectID string = "varity"

// DeploymentModeDev is a "dev" deployment
const DeploymentModeDev string = "dev"

// DeploymentModeProd is a "prod" deployment
const DeploymentModeProd string = "prod"

// BigqueryDatasetScraping refers to the scraping bigquery dataset
const BigqueryDatasetScraping string = "scraping"

// BigqueryTableRedditSubmissions refers to the reddit submissions v2 bigquery table
const BigqueryTableRedditSubmissions string = "reddit_submissions_v2"

// BigqueryTableRedditComments refers to the reddit comments v2 bigquery table
const BigqueryTableRedditComments string = "reddit_comments_v2"

// BigqueryTableTickerMentions refers to the ticker mentions v2 bigquery table
const BigqueryTableTickerMentions string = "ticker_mentions_v2"

// RedditSubmissions is just a general placeholder for "reddit-submissions"
const RedditSubmissions string = "reddit-submissions"

// RedditComments is just a general placeholder for "reddit-comments"
const RedditComments string = "reddit-comments"

// TickerMentions is just a general placeholder for "ticker-mentions"
const TickerMentions string = "ticker-mentions"

// KafkaPartitionsCount is the number of partitions used by each topic.
// This is defined in the confluent terraform files
const KafkaPartitionsCount int = 2

// ParentSourceRedditComment is a placeholder for the "reddit-comment" parent source
const ParentSourceRedditComment string = "reddit-comment"

// ParentSourceRedditSubmissionTitle is a placeholder for the "v" parent source
const ParentSourceRedditSubmissionTitle string = "reddit-submission-title"

// ParentSourceRedditSubmissionBody is a placeholder for the "reddit-submission-body" parent source
const ParentSourceRedditSubmissionBody string = "reddit-submission-body"
