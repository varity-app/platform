resource "google_pubsub_subscription" "etl_bigquery_to_influx_main" {
  name    = "etl-bigquery-to-influx-main-${var.deployment}"
  topic   = google_pubsub_topic.etl_bigquery_to_influx.name
  project = var.project

  ack_deadline_seconds = 60

  push_config {
    oidc_token {
      service_account_email = google_service_account.pubsub_invoker.email
    }

    push_endpoint = "${var.etl_bigquery_to_influx_url}/api/v1/etl/bigquery-to-influx"

    attributes = {
      x-goog-version = "v1"
    }
  }

  labels = {
    deployment = var.deployment
  }
}
