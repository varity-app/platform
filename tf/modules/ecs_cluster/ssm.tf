data "aws_ssm_parameter" "bootstrap" {
  name = format("varity-confluent-bootstrap%s", var.suffix)
}

data "aws_ssm_parameter" "confluent_key" {
  name = format("varity-confluent-key%s", var.suffix)
}

data "aws_ssm_parameter" "confluent_secret" {
  name = format("varity-confluent-secret%s", var.suffix)
}

data "aws_ssm_parameter" "reddit_client_id" {
  name = "varity-reddit-client-id"
}

data "aws_ssm_parameter" "reddit_client_secret" {
  name = "varity-reddit-client-secret"
}

data "aws_ssm_parameter" "reddit_username" {
  name = "varity-reddit-username"
}

data "aws_ssm_parameter" "reddit_password" {
  name = "varity-reddit-password"
}

data "aws_ssm_parameter" "reddit_user_agent" {
  name = "varity-reddit-user-agent"
}