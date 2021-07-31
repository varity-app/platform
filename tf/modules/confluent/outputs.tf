output "kafka_url" {
  value = local.bootstrap_servers
}

output "key" {
  value     = confluentcloud_api_key.svc.key
  sensitive = true
}

output "secret" {
  value     = confluentcloud_api_key.svc.secret
  sensitive = true
}