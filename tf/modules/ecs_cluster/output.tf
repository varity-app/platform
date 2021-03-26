output "cluster_id" {
    description = "ID of the ECS Cluster"
    value = aws_ecs_cluster.main.id
}