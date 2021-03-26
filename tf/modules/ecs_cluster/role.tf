data "aws_iam_policy_document" "ecs_agent" {
    statement {
        actions = [ "sts:AssumeRole" ]

        principals {
            type = "Service"
            identifiers = [ "ecs-tasks.amazonaws.com" ]
        }
    }
}

resource "aws_iam_policy" "cloudwatch" {
    name = "ECS-CloudWatch"
    description = "Policy for ECS to read and write to AWS CloudWatch"

    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogStreams"
            ],
            "Resource": [
                "arn:aws:logs:*:*:*"
            ]
        }
    ]
}
EOF
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

resource "aws_iam_policy_attachment" "read_parameters" {
    name = "ecs-agent-read-parameters-attachment"
    roles = [ aws_iam_role.ecs_agent.name ]
    policy_arn = var.secrets_policy_arn
}

resource "aws_iam_policy_attachment" "cloudwatch" {
    name = "ecs-agent-attachment"
    roles = [ aws_iam_role.ecs_agent.name ]
    policy_arn = aws_iam_policy.cloudwatch.arn
}