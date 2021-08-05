package common

// GCP
const GcpProjectId string = "varity"

// Deployment modes
const DeploymentModeDev string = "dev"
const DeploymentModeProd string = "prod"

// Bigquery
const BigqueryDatasetScraping string = "scraping"

const BigqueryTableRedditSubmissions string = "reddit_submissions_v2"
const BigqueryTableRedditComments string = "reddit_comments_v2"
const BigqueryTableTickerMentions string = "ticker_mentions_v2"

// General reddit
const RedditSubmissions string = "reddit-submissions"
const RedditComments string = "reddit-comments"

// General
const TickerMentions string = "ticker-mentions"

// Kafka
const KafkaPartitionsCount int = 2

// Data parent sources
const ParentSourceRedditComment string = "reddit-comment"
const ParentSourceRedditSubmissionTitle string = "reddit-submission-title"
const ParentSourceRedditSubmissionBody string = "reddit-submission-body"
