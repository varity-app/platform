package common

// GCP
const GCP_PROJECT_ID string = "varity"

// Deployment modes
const DEPLOYMENT_MODE_DEV string = "dev"
const DEPLOYMENT_MODE_PROD string = "prod"

// Pub/Sub subscription names
const SUBSCRIPTION_REDDIT_COMMENTS string = "reddit-comments-proc"
const SUBSCRIPTION_REDDIT_SUBSCRIPTIONS string = "reddit-submissions-proc"
const SUBSCRIPTION_TICKER_MENTIONS string = "ticker-mentions-proc"

// Pub/Sub topic names
const TOPIC_TICKER_MENTIONS string = "ticker-mentions"

// Bigquery
const BIGQUERY_DATASET_SCRAPING string = "scraping"

const BIGQUERY_TABLE_TICKER_MENTIONS string = "ticker_mentions_v2"

// General reddit
const REDDIT_SUBMISSIONS string = "reddit-submissions"
const REDDIT_COMMENTS string = "reddit-comments"

// General
const TICKER_MENTIONS string = "ticker-mentions"

// Kafka
const KAFKA_PARTITION_COUNT int = 2

// Data parent sources
const PARENT_SOURCE_REDDIT_COMMENT string = "reddit-comment"
const PARENT_SOURCE_REDDIT_SUBMISSION_TITLE string = "reddit-submission-title"
const PARENT_SOURCE_REDDIT_SUBMISSION_BODY string = "reddit-submission-body"
