data "aws_ami" "ecs_optimized" {
  owners = ["amazon", "aws-marketplace"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-2.0.20210316-x86_64-ebs"]
  }
}

resource "aws_launch_configuration" "ecs_launch_config" {
  name_prefix     = "varity-ecs"
  image_id = data.aws_ami.ecs_optimized.id

  associate_public_ip_address = true

  iam_instance_profile = aws_iam_instance_profile.ecs_agent.name
  security_groups      = [module.network.security_group]
  user_data            = format("#!/bin/bash\necho ECS_CLUSTER=%s >> /etc/ecs/ecs.config", var.ecs_cluster_name)

  instance_type = "t3a.micro"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "ecs" {
  name                 = "varity-ecs-asg"
  vpc_zone_identifier  = [module.network.subnet_id]
  launch_configuration = aws_launch_configuration.ecs_launch_config.name

  desired_capacity          = 2
  min_size                  = 0
  max_size                  = 5
  health_check_grace_period = 300
  health_check_type         = "EC2"
  force_delete              = true
}
