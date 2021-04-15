resource "google_bigquery_table" "reddit_submissions" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "reddit_submissions"

  time_partitioning {
    type  = "DAY"
    field = "created_utc"
  }

  labels = {
    deployment = var.deployment
  }

  schema = <<EOF
[
  {
    "name": "submission_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "subreddit",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "title",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "created_utc",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "name",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "selftext",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "author",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "is_original_content",
    "type": "BOOL",
    "mode": "REQUIRED"
  },
  {
    "name": "is_text",
    "type": "BOOL",
    "mode": "REQUIRED"
  },
  {
    "name": "nsfw",
    "type": "BOOL",
    "mode": "REQUIRED"
  },
  {
    "name": "num_comments",
    "type": "INT64",
    "mode": "NULLABLE"
  },
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "upvotes",
    "type": "INT64",
    "mode": "NULLABLE"
  },
  {
    "name": "upvote_ratio",
    "type": "FLOAT64",
    "mode": "NULLABLE"
  },
  {
    "name": "url",
    "type": "STRING",
    "mode": "REQUIRED"
  }
]
EOF

}

resource "google_bigquery_table" "reddit_comments" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "reddit_comments"

  time_partitioning {
    type  = "DAY"
    field = "created_utc"
  }

  labels = {
    deployment = var.deployment
  }

  schema = <<EOF
[
  {
    "name": "comment_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "submission_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "subreddit",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "author",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "created_utc",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "body",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "upvotes",
    "type": "INT64",
    "mode": "NULLABLE"
  }
]
EOF

}

resource "google_bigquery_table" "scraped_posts" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "scraped_posts"

  time_partitioning {
    type  = "DAY"
    field = "timestamp"
  }

  labels = {
    deployment = var.deployment
  }

  schema = <<EOF
[
  {
    "name": "data_source",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "parent_source",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "parent_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "text",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  }
]
EOF

}

resource "google_bigquery_table" "ticker_mentions" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "ticker_mentions"

  time_partitioning {
    type  = "DAY"
    field = "timestamp"
  }

  labels = {
    deployment = var.deployment
  }

  schema = <<EOF
[
  {
    "name": "ticker",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "data_source",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "parent_source",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "parent_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "mention_type",
    "type": "STRING",
    "mode": "REQUIRED"
  }
]
EOF

}
