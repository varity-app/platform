terraform {
  required_providers {
    kafka = {
      source = "Mongey/kafka"
      version = "0.2.11"
    }
  }
}

provider "kafka" {
    bootstrap_servers = [var.bootstrap_servers]
    tls_enabled = true
    sasl_username = var.confluent_key
    sasl_password = var.confluent_secret
}