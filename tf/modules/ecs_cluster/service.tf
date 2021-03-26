resource "aws_ecs_service" "scrape_reddit_comments" {
  name            = "scrape-reddit-comments"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.comments_scraper.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  network_configuration {
    subnets          = [module.network.subnet_id]
    security_groups  = [module.network.security_group]
    assign_public_ip = true
  }
}

resource "aws_ecs_service" "scrape_reddit_submissions" {
  name            = "scrape-reddit-submissions"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.submissions_scraper.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  network_configuration {
    subnets          = [module.network.subnet_id]
    security_groups  = [module.network.security_group]
    assign_public_ip = true
  }
}

resource "aws_ecs_service" "faust" {
  name            = "faust"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.faust.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  network_configuration {
    subnets          = [module.network.subnet_id]
    security_groups  = [module.network.security_group]
    assign_public_ip = true
  }
}