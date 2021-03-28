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

variable "asg_name" {
  description = "Name of the EC2 autoscaling group used for ECS"
  type        = string
}

variable "ecs_role_name" {
  description = "Name of the ECS iam execution role to be created"
  type        = string
}

variable "ecs_instance_profile_name" {
  description = "Name of the ECS iam instance profile to be created"
  type        = string
}

variable "ecs_cloudwatch_policy_name" {
  description = "Name of the IAM policy that lets the ECS execution role write to CloudWatch"
  type        = string
}
