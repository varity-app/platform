variable "ecs_cluster_name" {
    description = "Name of the ECS Cluster"
    type = string
    default = "varity-dev"
}

variable "security_group_name" {
    description = "Name of the VPC security group that ECS will use"
    type = string
    default = "varity-ecs-dev"
}