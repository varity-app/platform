resource "google_bigquery_table" "reddit_submissions" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "reddit_submissions_v2"

  time_partitioning {
    type  = "DAY"
    field = "timestamp"
  }

  labels = {
    deployment = var.deployment
  }


  lifecycle {
    prevent_destroy = true
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
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "body",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "author",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "author_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "is_self",
    "type": "BOOL",
    "mode": "REQUIRED"
  },
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "url",
    "type": "STRING",
    "mode": "NULLABLE"
  }
]
EOF

}

resource "google_bigquery_table" "reddit_comments" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "reddit_comments_v2"

  time_partitioning {
    type  = "DAY"
    field = "timestamp"
  }

  labels = {
    deployment = var.deployment
  }

  lifecycle {
    prevent_destroy = true
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
    "name": "author_id",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "timestamp",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "body",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "REQUIRED"
  }
]
EOF

}

resource "google_bigquery_table" "ticker_mentions" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "ticker_mentions_v2"

  time_partitioning {
    type  = "DAY"
    field = "timestamp"
  }

  labels = {
    deployment = var.deployment
  }

  lifecycle {
    prevent_destroy = true
  }

  schema = <<EOF
[
  {
    "name": "symbol",
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
    "name": "symbol_counts",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "short_name_counts",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "word_count",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "question_mark_count",
    "type": "INT64",
    "mode": "REQUIRED"
  }
]
EOF

}

resource "google_bigquery_table" "eod_prices" {
  project    = var.project
  dataset_id = google_bigquery_dataset.scraping.dataset_id
  table_id   = "eod_prices"

  time_partitioning {
    type  = "DAY"
    field = "date"
  }

  labels = {
    deployment = var.deployment
  }

  lifecycle {
    prevent_destroy = true
  }

  schema = <<EOF
[
  {
    "name": "symbol",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "date",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  },
  {
    "name": "open",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "high",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "low",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "close",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "volume",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "adj_open",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "adj_high",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "adj_low",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "adj_close",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "adj_volume",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "dividend",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "split",
    "type": "FLOAT",
    "mode": "REQUIRED"
  },
  {
    "name": "scraped_on",
    "type": "TIMESTAMP",
    "mode": "REQUIRED"
  }
]
EOF

}
