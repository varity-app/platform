resource "google_pubsub_topic" "reddit_submissions" {
    name = format("reddit-submissions-%s", var.deployment)
    project = var.project

    labels = {
        deployment = var.deployment
    }
}

resource "google_pubsub_topic" "reddit_comments" {
    name = format("reddit-comments-%s", var.deployment)
    project = var.project

    labels = {
        deployment = var.deployment
    }
}