resource "aws_ecs_task_definition" "test" {
    family = "test-definition"
    network_mode = "awsvpc"
    requires_compatibilities = [ "FARGATE" ]
    cpu = 256
    memory = 512
    task_role_arn = aws_iam_role.ecs_agent.arn
    execution_role_arn = aws_iam_role.ecs_agent.arn

    container_definitions = jsonencode([
        {
            name = "test-image"
            image = "dockerbogo/docker-nginx-hello-world"
            essential = true
            cpu = 256
            memory = 512
            portMappings = [
                {
                    containerPort = 80
                    hostPort = 80
                }
            ]
        }
    ])
}