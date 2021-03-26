data "aws_ssm_parameter" "bootstrap" {
    name = "varity-confluent-bootstrap"
}

data "aws_ssm_parameter" "confluent_key" {
    name = "varity-confluent-key"
}

data "aws_ssm_parameter" "confluent_secret" {
    name = "varity-confluent-secret"
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