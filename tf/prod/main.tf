terraform {
  backend "remote" {
    organization = "Varity"

    workspaces {
      name = "varity-prod"
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

module "kafka" {
  source = "../modules/kafka"

  bootstrap_servers = var.bootstrap_servers
  confluent_key     = var.confluent_key_prod
  confluent_secret  = var.confluent_secret_prod
  num_partitions    = var.num_partitions
}

module "ecs_cluster" {
  source = "../modules/ecs_cluster"

  ecs_cluster_name           = var.ecs_cluster_name
  security_group_name        = var.security_group_name
  secrets_policy_arn         = var.secrets_policy_arn
  asg_name                   = var.asg_name
  ecs_role_name              = var.ecs_role_name
  ecs_instance_profile_name  = var.ecs_instance_profile_name
  ecs_cloudwatch_policy_name = var.ecs_cloudwatch_policy_name
  suffix                     = var.suffix
  submissions_table_name     = var.submissions_table_name
  comments_table_name        = var.comments_table_name
}

module "dynamo" {
  source                 = "../modules/dynamo"
  submissions_table_name = var.submissions_table_name
  comments_table_name    = var.comments_table_name
}
