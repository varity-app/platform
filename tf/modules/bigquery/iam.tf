resource "google_service_account" "metabase" {
  account_id   = "varity-metabase-svc-${var.deployment}"
  display_name = "Varity Metabase Service Account (${var.deployment})"
  project      = var.project
}

resource "google_project_iam_member" "data_viewer" {
  project = var.project
  role    = "roles/bigquery.dataViewer"
  member  = "serviceAccount:${google_service_account.metabase.email}"
}

resource "google_project_iam_member" "metadata_viewer" {
  project = var.project
  role    = "roles/bigquery.metadataViewer"
  member  = "serviceAccount:${google_service_account.metabase.email}"
}

resource "google_project_iam_member" "job_user" {
  project = var.project
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.metabase.email}"
}

resource "google_service_account_key" "metabase_key" {
  service_account_id = google_service_account.metabase.name
}
