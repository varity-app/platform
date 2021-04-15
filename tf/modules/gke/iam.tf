resource "google_project_iam_binding" "pubsub_publish" {
  project = var.project

  role = "roles/pubsub.publisher"

  members = [
    "serviceAccount:${google_service_account.default.email}"
  ]
}

resource "google_project_iam_binding" "pubsub_subscribe" {
  project = var.project

  role = "roles/pubsub.subscriber"

  members = [
    "serviceAccount:${google_service_account.default.email}"
  ]
}

resource "google_project_iam_binding" "firestore" {
  project = var.project

  role = "roles/datastore.owner"

  members = [
    "serviceAccount:${google_service_account.default.email}"
  ]
}

resource "google_project_iam_binding" "bigquery" {
  project = var.project

  role = "roles/bigquery.dataEditor"

  members = [
    "serviceAccount:${google_service_account.default.email}"
  ]
}

resource "google_service_account" "default" {
  account_id   = "varity-gke-svc-${var.deployment}"
  display_name = "Varity GKE Service Account (${var.deployment})"
  project      = var.project
}
