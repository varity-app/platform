resource "google_bigquery_dataset" "scraping" {
  project                    = var.project
  dataset_id                 = "scraping_${var.deployment}"
  friendly_name              = "Varity Scraping (${var.deployment})"
  description                = "General warehouse for varity social media scraping"
  location                   = "US"
  delete_contents_on_destroy = true

  labels = {
    deployment = var.deployment
  }
}

resource "google_bigquery_dataset" "scraping_dbt" {
  project                    = var.project
  dataset_id                 = "dbt_scraping_${var.deployment}"
  friendly_name              = "Varity Scraping DBT (${var.deployment})"
  description                = "DBT models for varity social media scraping"
  location                   = "US"
  delete_contents_on_destroy = true

  labels = {
    deployment = var.deployment
  }
}
