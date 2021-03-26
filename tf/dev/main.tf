terraform {
    required_providers {
      aws = {
          source = "hashicorp/aws"
          version = "~> 3.0"
      }
    }
}

provider "aws" {
    region = "us-east-2"
}

module "ecs_cluster" {
    source = "../modules/ecs_cluster"
}