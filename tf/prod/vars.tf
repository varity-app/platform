variable "ecs_cluster_name" {
  description = "Name of the ECS Cluster"
  type        = string
  default     = "varity-prod"
}

variable "security_group_name" {
  description = "Name of the VPC security group that ECS will use"
  type        = string
  default     = "varity-ecs-prod"
}

variable "secrets_policy_arn" {
  description = "ARN of the IAM policy that permits reading of AWS SSM secrets"
  type        = string
  default     = "arn:aws:iam::178852309825:policy/readVaritySecrets"
}

variable "asg_name" {
  description = "Name of the EC2 autoscaling group used for ECS"
  type        = string
  default     = "varity-ecs-asg-prod"
}

variable "ecs_role_name" {
  description = "Name of the ECS execution role to be created"
  type        = string
  default     = "varityECSExecutionRoleProd"
}

variable "ecs_instance_profile_name" {
  description = "Name of the ECS iam instance profile to be created"
  type        = string
  default     = "ecs-agent-prod"
}

variable "ecs_cloudwatch_policy_name" {
  description = "Name of the IAM policy that lets the ECS execution role write to CloudWatch"
  type        = string
  default     = "ECS-CloudWatch-Prod"
}

variable "suffix" {
  description = "Suffix to append to various definitions to seperate deployments (e.g `-dev`, `-prod`)"
  type        = string
  default     = "-prod"
}

variable "submissions_table_name" {
  description = "Name of the reddit submissions DynamoDB table name"
  type        = string
  default     = "reddit-submissions-prod"
}

variable "comments_table_name" {
  description = "Name of the reddit submissions DynamoDB table name"
  type        = string
  default     = "reddit-comments-prod"
}

variable "bootstrap_servers" {
  description = "URL of the Confluent Cloud boostrap servers"
  type        = string
  default     = "pkc-ep9mm.us-east-2.aws.confluent.cloud:9092"
}

variable "confluent_key_prod" {
  description = "SASL username to login to the production confluent cluster"
  type        = string
}

variable "confluent_secret_prod" {
  description = "SASL password to login to the production confluent cluster"
  type        = string
}

variable "num_partitions" {
  description = "Number of partitions to use for each Kafka topic"
  type        = number
  default     = 2
}