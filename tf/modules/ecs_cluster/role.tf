data "aws_iam_policy_document" "ecs_agent" {
    statement {
        actions = [ "sts:AssumeRole" ]

        principals {
            type = "Service"
            identifiers = [ "ecs-tasks.amazonaws.com" ]
        }
    }
}

resource "aws_iam_role" "ecs_agent" {
    name = "varityECSExecutionRole"

    assume_role_policy = data.aws_iam_policy_document.ecs_agent.json
}

resource "aws_iam_policy_attachment" "ecs_agent" {
    name = "ecs-agent-attachment"
    roles = [ aws_iam_role.ecs_agent.name ]
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}