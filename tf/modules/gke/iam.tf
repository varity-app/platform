resource "google_project_iam_member" "pubsub_publish" {
  project = var.project
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "pubsub_subscribe" {
  project = var.project
  role = "roles/pubsub.subscriber"
  member = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "firestore" {
  project = var.project
  role = "roles/datastore.owner"
  member = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "bigquery" {
  project = var.project
  role = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.default.email}"
}

resource "google_service_account" "default" {
  account_id   = "varity-gke-svc-${var.deployment}"
  display_name = "Varity GKE Service Account (${var.deployment})"
  project      = var.project
}
