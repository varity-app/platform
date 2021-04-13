terraform {
  backend "remote" {
    organization = "Varity"

    workspaces {
      name = "varity-prod"
    }
  }

  required_providers {
    aws = {
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
