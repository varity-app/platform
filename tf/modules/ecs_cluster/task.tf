resource "aws_ecs_task_definition" "comments_scraper" {
  family                   = "comments-scraper-definition"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  cpu    = 256
  memory = 512

  task_role_arn      = aws_iam_role.ecs_agent.arn
  execution_role_arn = aws_iam_role.ecs_agent.arn

  container_definitions = jsonencode([
    {
      name      = "scraper"
      image     = "cgundlach13/reddit-scraper:0.3.0"
      essential = true

      cpu    = 256
      memory = 512

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/scrapers/reddit-comments"
          awslogs-stream-prefix = "ecs"
          awslogs-region        = "us-east-2"
          awslogs-create-group  = "true"
        }
      }

      environment = [
        {
          name  = "MODE"
          value = "comments"
        },
        {
          name  = "SUBREDDITS"
          value = "wallstreetbets,smallstreetbets,stocks,valueinvesting,securityanalysis,investing"
        },
        {
          name  = "SLEEP"
          value = "60"
        }
      ]

      secrets = [
        {
          name      = "REDDIT_USERNAME"
          valueFrom = data.aws_ssm_parameter.reddit_username.arn
        },
        {
          name      = "REDDIT_PASSWORD"
          valueFrom = data.aws_ssm_parameter.reddit_password.arn
        },
        {
          name      = "REDDIT_CLIENT_ID"
          valueFrom = data.aws_ssm_parameter.reddit_client_id.arn
        },
        {
          name      = "REDDIT_CLIENT_SECRET"
          valueFrom = data.aws_ssm_parameter.reddit_client_secret.arn
        },
        {
          name      = "REDDIT_USER_AGENT"
          valueFrom = data.aws_ssm_parameter.reddit_user_agent.arn
        },
        {
          name      = "SASL_USERNAME"
          valueFrom = data.aws_ssm_parameter.confluent_key.arn
        },
        {
          name      = "SASL_PASSWORD"
          valueFrom = data.aws_ssm_parameter.confluent_secret.arn
        },
        {
          name      = "BOOTSTRAP_SERVERS"
          valueFrom = data.aws_ssm_parameter.bootstrap.arn
        }
      ]
    }
  ])
}

resource "aws_ecs_task_definition" "submissions_scraper" {
  family                   = "submissions-scraper-definition"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  cpu    = 256
  memory = 512

  task_role_arn      = aws_iam_role.ecs_agent.arn
  execution_role_arn = aws_iam_role.ecs_agent.arn

  container_definitions = jsonencode([
    {
      name      = "scraper"
      image     = "cgundlach13/reddit-scraper:0.3.0"
      essential = true

      cpu    = 256
      memory = 512

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/scrapers/reddit-submissions"
          awslogs-stream-prefix = "ecs"
          awslogs-region        = "us-east-2"
          awslogs-create-group  = "true"
        }
      }

      environment = [
        {
          name  = "MODE"
          value = "submissions"
        },
        {
          name  = "SUBREDDITS"
          value = "wallstreetbets,smallstreetbets,stocks,valueinvesting,securityanalysis,investing"
        },
        {
          name  = "SLEEP"
          value = "120"
        }
      ]

      secrets = [
        {
          name      = "REDDIT_USERNAME"
          valueFrom = data.aws_ssm_parameter.reddit_username.arn
        },
        {
          name      = "REDDIT_PASSWORD"
          valueFrom = data.aws_ssm_parameter.reddit_password.arn
        },
        {
          name      = "REDDIT_CLIENT_ID"
          valueFrom = data.aws_ssm_parameter.reddit_client_id.arn
        },
        {
          name      = "REDDIT_CLIENT_SECRET"
          valueFrom = data.aws_ssm_parameter.reddit_client_secret.arn
        },
        {
          name      = "REDDIT_USER_AGENT"
          valueFrom = data.aws_ssm_parameter.reddit_user_agent.arn
        },
        {
          name      = "SASL_USERNAME"
          valueFrom = data.aws_ssm_parameter.confluent_key.arn
        },
        {
          name      = "SASL_PASSWORD"
          valueFrom = data.aws_ssm_parameter.confluent_secret.arn
        },
        {
          name      = "BOOTSTRAP_SERVERS"
          valueFrom = data.aws_ssm_parameter.bootstrap.arn
        }
      ]
    }
  ])
}

resource "aws_ecs_task_definition" "faust" {
  family                   = "faust-definition"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  cpu    = 256
  memory = 512

  task_role_arn      = aws_iam_role.ecs_agent.arn
  execution_role_arn = aws_iam_role.ecs_agent.arn

  container_definitions = jsonencode([
    {
      name      = "scraper"
      image     = "cgundlach13/faust-processor:0.3.0"
      essential = true

      cpu    = 256
      memory = 512

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/faust"
          awslogs-stream-prefix = "ecs"
          awslogs-region        = "us-east-2"
          awslogs-create-group  = "true"
        }
      }

      secrets = [
        {
          name      = "SASL_USERNAME"
          valueFrom = data.aws_ssm_parameter.confluent_key.arn
        },
        {
          name      = "SASL_PASSWORD"
          valueFrom = data.aws_ssm_parameter.confluent_secret.arn
        },
        {
          name      = "BOOTSTRAP_SERVERS"
          valueFrom = data.aws_ssm_parameter.bootstrap.arn
        }
      ]

      mountPoints = [
        {
          sourceVolume  = "faust-vol"
          containerPath = "/tmp/faust"
          readOnly      = false
        }
      ]
    }
  ])

  volume {
    name = "faust-vol"
  }
}