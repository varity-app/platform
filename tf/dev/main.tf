terraform {
  backend "remote" {
    organization = "Varity"

    workspaces {
      name = "varity-dev"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}

module "ecs_cluster" {
  source = "../modules/ecs_cluster"

  ecs_cluster_name    = var.ecs_cluster_name
  security_group_name = var.security_group_name
  secrets_policy_arn  = var.secrets_policy_arn
}