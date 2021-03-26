variable "ecs_cluster_name" {
  description = "Name of the ECS Cluster"
  type        = string
}

variable "security_group_name" {
  description = "Name of the VPC security group that ECS will use"
  type        = string
}

variable "secrets_policy_arn" {
  description = "ARN of the IAM policy that permits reading of AWS SSM secrets"
  type        = string
}