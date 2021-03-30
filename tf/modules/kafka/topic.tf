resource "kafka_topic" "reddit_submissions" {
    name = "reddit-submissions"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "reddit_comments" {
    name = "reddit-comments"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "ticker_mentions" {
    name = "ticker-mentions"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "scraped_posts" {
    name = "scraped-posts"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "post_sentiment" {
    name = "post-sentiment"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "faust" {
    name = "varity-faust-app-__assignor-__leader"
    partitions = var.num_partitions
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}

resource "kafka_topic" "logs" {
    name = "logs"
    partitions = 1
    replication_factor = 3

    config = {
        "retention.ms" = "604800000"
        "cleanup.policy" = "delete"
    }
}