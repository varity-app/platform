resource "google_pubsub_subscription" "submissions_proc" {
  name    = "reddit-submissions-proc-${var.deployment}"
  topic   = google_pubsub_topic.reddit_submissions.name
  project = var.project

  labels = {
    deployment = var.deployment
  }
}

resource "google_pubsub_subscription" "comments_proc" {
  name    = "reddit-comments-proc-${var.deployment}"
  topic   = google_pubsub_topic.reddit_comments.name
  project = var.project

  labels = {
    deployment = var.deployment
  }
}

resource "google_pubsub_subscription" "ticker_mentions_proc" {
  name    = "ticker-mentions-proc-${var.deployment}"
  topic   = google_pubsub_topic.ticker_mentions.name
  project = var.project

  labels = {
    deployment = var.deployment
  }
}
