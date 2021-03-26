resource "aws_ecs_service" "test" {
    name = "test"
    cluster = aws_ecs_cluster.main.id
    task_definition = aws_ecs_task_definition.test.arn
    desired_count = 1
    launch_type = "FARGATE"

    deployment_minimum_healthy_percent = 100
    deployment_maximum_percent = 200

    network_configuration {
      subnets = [ module.network.subnet_id ]
      security_groups = [ module.network.security_group ]
      assign_public_ip = true
    }
}