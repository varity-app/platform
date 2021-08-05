terraform {
  required_providers {
    confluentcloud = {
      source  = "Mongey/confluentcloud"
      version = "~> 0.0.11"
    }
    kafka = {
      source  = "Mongey/kafka"
      version = "0.2.11"
    }
  }
}

resource "confluentcloud_environment" "environment" {
  name = var.deployment
}

resource "confluentcloud_kafka_cluster" "scraping" {
  name             = "scraping"
  service_provider = "gcp"
  region           = var.region
  availability     = "LOW"
  environment_id   = confluentcloud_environment.environment.id
  deployment = {
    sku = "BASIC"
  }
  network_egress  = 100
  network_ingress = 100
  storage         = 5000
}

resource "confluentcloud_api_key" "svc" {

  cluster_id     = confluentcloud_kafka_cluster.scraping.id
  environment_id = confluentcloud_environment.environment.id
}

# Format bootstrap servers URL for the kafka provider
locals {
  bootstrap_servers = [replace(confluentcloud_kafka_cluster.scraping.bootstrap_servers, "SASL_SSL://", "")]
}
