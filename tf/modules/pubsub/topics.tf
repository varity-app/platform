resource "google_pubsub_topic" "etl_bigquery_to_influx" {
  name    = "etl-bigquery-to-influx-${var.deployment}"
  project = var.project

  labels = {
    deployment = var.deployment
  }
}
