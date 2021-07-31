# Secret containing the bootstrap server url
resource "google_secret_manager_secret" "kafka_url" {
  secret_id = "kafka-scraping-url-${var.deployment}"

  replication {
    automatic = true
  }
}


resource "google_secret_manager_secret_version" "kafka_url" {
  secret = google_secret_manager_secret.kafka_url.id

  secret_data = local.bootstrap_servers[0]
}

# Secret containing the auth key
resource "google_secret_manager_secret" "kafka_key" {
  secret_id = "kafka-scraping-key-${var.deployment}"

  replication {
    automatic = true
  }
}


resource "google_secret_manager_secret_version" "kafka_key" {
  secret = google_secret_manager_secret.kafka_key.id

  secret_data = confluentcloud_api_key.svc.key
}

# Secret containing the auth secret
resource "google_secret_manager_secret" "kafka_secret" {
  secret_id = "kafka-scraping-secret-${var.deployment}"

  replication {
    automatic = true
  }
}


resource "google_secret_manager_secret_version" "kafka_secret" {
  secret = google_secret_manager_secret.kafka_secret.id

  secret_data = confluentcloud_api_key.svc.secret
}