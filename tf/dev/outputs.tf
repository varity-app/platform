output "kafka_url" {
  value = module.confluent.kafka_url
}

output "kafka_key" {
  value     = module.confluent.key
  sensitive = true
}

output "kafka_secret" {
  value     = module.confluent.secret
  sensitive = true
}