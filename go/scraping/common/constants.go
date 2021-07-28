package common

// GCP
const GCP_PROJECT_ID string = "varity"

// Deployment modes
const DEPLOYMENT_MODE_DEV = "dev"
const DEPLOYMENT_MODE_PROD = "prod"

// Pub/Sub subscription names
const SUBSCRIPTION_REDDIT_COMMENTS = "reddit-comments-proc"
const SUBSCRIPTION_REDDIT_SUBSCRIPTIONS = "reddit-submissions-proc"

// Pub/Sub topic names
const TOPIC_TICKER_MENTIONS = "ticker-mentions"

// Data parent sources
const PARENT_SOURCE_REDDIT_COMMENT string = "reddit-comment"
const PARENT_SOURCE_REDDIT_SUBMISSION_TITLE string = "reddit-submission-title"
const PARENT_SOURCE_REDDIT_SUBMISSION_BODY string = "reddit-submission-body"
