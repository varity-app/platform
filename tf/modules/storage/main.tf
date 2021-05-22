resource "google_storage_bucket" "prefect" {
  name          = "varity-prefect-${var.deployment}"
  location      = "US"
  force_destroy = true
}