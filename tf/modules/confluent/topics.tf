provider "kafka" {
  bootstrap_servers = local.bootstrap_servers

  tls_enabled    = true
  sasl_username  = confluentcloud_api_key.svc.key
  sasl_password  = confluentcloud_api_key.svc.secret
  sasl_mechanism = "plain"
  timeout        = 10
}

resource "kafka_topic" "reddit_submissions" {
  name               = "reddit_submissions"
  replication_factor = 3
  partitions         = 2
  config = {
    "cleanup.policy" = "delete"
    "retention.ms" = "604800000"
  }
}

resource "kafka_topic" "reddit_comments" {
  name               = "reddit_comments"
  replication_factor = 3
  partitions         = 2
  config = {
    "cleanup.policy" = "delete"
    "retention.ms" = "604800000"
  }
}

resource "kafka_topic" "ticker_mentions" {
  name               = "ticker_mentions"
  replication_factor = 3
  partitions         = 2
  config = {
    "cleanup.policy" = "delete"
    "retention.ms" = "604800000"
  }
}