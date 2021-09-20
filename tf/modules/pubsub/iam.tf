// Service account used by Pub/Sub to invoke Cloud Run
resource "google_service_account" "pubsub_invoker" {
  account_id   = "pubsub-invoker-svc-${var.deployment}"
  display_name = "Pub/Sub Cloud Run Invoker"
}

// Give the roles/run.invoker role to the service account
resource "google_project_iam_member" "pubsub_invoker_cloud_run_invoke" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.pubsub_invoker.email}"
}