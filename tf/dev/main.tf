terraform {
  backend "remote" {
    organization = "Varity"

    workspaces {
      name = "varity-dev"
    }
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.64.0"
    }
  }
}

provider "google" {
  project = "varity"
  region  = "us-east1"
  zone    = "us-east1-c"
}

module "pubsub" {
  source = "../modules/pubsub"

  deployment = var.deployment
}

# module "gke" {
#   source = "../modules/gke"

#   deployment = var.deployment
#   release    = var.release
# }

module "biquery" {
  source = "../modules/bigquery"

  deployment = var.deployment
}

module "cloud_sql" {
  source = "../modules/cloud_sql"

  deployment = var.deployment
}

module "confluent" {
  source = "../modules/confluent"

  deployment = var.deployment
}

module "cloud_run" {
  source = "../modules/cloud_run"

  deployment                = var.deployment
  release                   = var.release
  cloud_sql_connection_name = module.cloud_sql.connection_name

  kafka_url_secret_id    = module.confluent.kafka_url_secret_id
  kafka_key_secret_id    = module.confluent.kafka_key_secret_id
  kafka_secret_secret_id = module.confluent.kafka_secret_secret_id
}