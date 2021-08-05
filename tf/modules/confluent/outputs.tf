output "kafka_url" {
  value = local.bootstrap_servers
}

output "kafka_url_secret_id" {
  value     = google_secret_manager_secret.kafka_url.secret_id
  sensitive = true
}

output "key" {
  value     = confluentcloud_api_key.svc.key
  sensitive = true
}

output "kafka_key_secret_id" {
  value     = google_secret_manager_secret.kafka_key.secret_id
  sensitive = true
}

output "secret" {
  value     = confluentcloud_api_key.svc.secret
  sensitive = true
}

output "kafka_secret_secret_id" {
  value     = google_secret_manager_secret.kafka_secret.secret_id
  sensitive = true
}