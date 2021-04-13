resource "google_pubsub_topic" "reddit_submissions" {
  name    = "reddit-submissions-${var.deployment}"
  project = var.project

  labels = {
    deployment = var.deployment
  }
}

resource "google_pubsub_topic" "reddit_comments" {
  name    = "reddit-comments-${var.deployment}"
  project = var.project

  labels = {
    deployment = var.deployment
  }
}
