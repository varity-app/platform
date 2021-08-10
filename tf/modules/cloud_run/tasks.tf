resource "google_cloud_tasks_queue" "historical_reddit_queue" {
  name = "reddit-historical-${var.deployment}"
  location = var.region
  project  = var.project

  rate_limits {
    max_concurrent_dispatches = 8
    max_dispatches_per_second = 1
  }

  retry_config {
    max_attempts = 10
    max_retry_duration = "4s"
    max_backoff = "1000s"
    min_backoff = "10s"
  }
}

resource "google_project_iam_member" "tasks_svc_invoker" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.tasks_svc.email}"
}

resource "google_service_account" "tasks_svc" {
  account_id   = "varity-tasks-svc-${var.deployment}"
  display_name = "Varity Cloud Tasks Service Account (${var.deployment})"
  project      = var.project
}