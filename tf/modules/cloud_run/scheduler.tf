
// Service account used by Cloud Scheduler
resource "google_service_account" "scheduler_svc" {
  account_id   = "cloud-scheduler-svc-${var.deployment}"
  display_name = "Cloud Scheduler Service Account"
}

// Give the roles/run.invoker role to the service account
resource "google_project_iam_member" "scheduler_cloud_run_invoke" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.scheduler_svc.email}"
}

/* Scraping Jobs */
resource "google_cloud_scheduler_job" "scrape_reddit_submissions_wallstreetbets" {
  name             = "scrape-reddit-submissions-wallstreetbets-${var.deployment}"
  description      = "Scrape reddit submissions from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/submissions/wallstreetbets"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

resource "google_cloud_scheduler_job" "scrape_reddit_submissions_smallstreetbets" {
  name             = "scrape-reddit-submissions-smallstreetbets-${var.deployment}"
  description      = "Scrape reddit submissions from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/submissions/smallstreetbets"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

resource "google_cloud_scheduler_job" "scrape_reddit_submissions_stocks" {
  name             = "scrape-reddit-submissions-stocks-${var.deployment}"
  description      = "Scrape reddit submissions from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/submissions/stocks"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

resource "google_cloud_scheduler_job" "scrape_reddit_comments_wallstreetbets" {
  name             = "scrape-reddit-comments-wallstreetbets-${var.deployment}"
  description      = "Scrape reddit comments from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/comments/wallstreetbets"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}


resource "google_cloud_scheduler_job" "scrape_reddit_comments_smallstreetbets" {
  name             = "scrape-reddit-comments-smallstreetbets-${var.deployment}"
  description      = "Scrape reddit comments from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/comments/smallstreetbets"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

resource "google_cloud_scheduler_job" "scrape_reddit_comments_stocks" {
  name             = "scrape-reddit-comments-stocks-${var.deployment}"
  description      = "Scrape reddit comments from a specified subreddit"
  schedule         = "*/1 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.scrape_reddit.status[0].url}/scraping/reddit/comments/stocks"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

/* Processing Jobs */
resource "google_cloud_scheduler_job" "proc_reddit_comments_tickers" {
  name             = "proc-reddit-comments-tickers-${var.deployment}"
  description      = "Process reddit comments for tickers"
  schedule         = "*/5 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.proc.status[0].url}/scraping/proc/reddit/comments/extract"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}


resource "google_cloud_scheduler_job" "proc_reddit_comments_sink" {
  name             = "proc-reddit-comments-sink-${var.deployment}"
  description      = "Sink reddit comments to bigquery"
  schedule         = "*/5 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.proc.status[0].url}/scraping/proc/reddit/comments/sink"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}


resource "google_cloud_scheduler_job" "proc_reddit_submissions_tickers" {
  name             = "proc-reddit-submissions-tickers-${var.deployment}"
  description      = "Process reddit submissions for tickers"
  schedule         = "*/5 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.proc.status[0].url}/scraping/proc/reddit/submissions/extract"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}


resource "google_cloud_scheduler_job" "proc_reddit_submissions_sink" {
  name             = "proc-reddit-submissions-sink-${var.deployment}"
  description      = "Sink reddit submissions to bigquery"
  schedule         = "*/5 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.proc.status[0].url}/scraping/proc/reddit/submissions/sink"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}

resource "google_cloud_scheduler_job" "proc_ticker_mentions_sink" {
  name             = "proc-ticker-mentions-sink-${var.deployment}"
  description      = "Sink ticker mentions to bigquery"
  schedule         = "*/5 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.proc.status[0].url}/scraping/proc/tickerMentions/sink"

    oidc_token {
      service_account_email = google_service_account.scheduler_svc.email
    }
  }
}