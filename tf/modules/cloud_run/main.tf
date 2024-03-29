module "secrets" {
  source = "../secrets"
}

resource "google_cloud_run_service" "proc" {
  name     = "proc-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 8
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/proc:${var.release}"

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }

        env {
          name  = "POSTGRES_ADDRESS"
          value = "/cloudsql/${var.cloud_sql_connection_name}/.s.PGSQL.5432"
        }

        env {
          name  = "POSTGRES_NETWORK"
          value = "unix"
        }

        env {
          name  = "POSTGRES_DB"
          value = "finance"
        }

        env {
          name = "POSTGRES_USERNAME"
          value_from {
            secret_key_ref {
              name = module.secrets.postgres_username_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_PASSWORD"
          value_from {
            secret_key_ref {
              name = module.secrets.postgres_password_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_BOOTSTRAP_SERVERS"
          value_from {
            secret_key_ref {
              name = var.kafka_url_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_KEY"
          value_from {
            secret_key_ref {
              name = var.kafka_key_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_SECRET"
          value_from {
            secret_key_ref {
              name = var.kafka_secret_secret_id
              key  = "latest"
            }
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "10"
        "run.googleapis.com/cloudsql-instances" = var.cloud_sql_connection_name
      }
    }
  }

  metadata {
    annotations = {
      generated-by                      = "magic-modules"
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}

resource "google_cloud_run_service" "scrape_reddit" {
  name     = "scrape-reddit-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 8
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/reddit-scraper:${var.release}"

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }

        env {
          name = "REDDIT_CLIENT_ID"
          value_from {
            secret_key_ref {
              name = module.secrets.reddit_client_id_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDDIT_CLIENT_SECRET"
          value_from {
            secret_key_ref {
              name = module.secrets.reddit_client_secret_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDDIT_USERNAME"
          value_from {
            secret_key_ref {
              name = module.secrets.reddit_username_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDDIT_PASSWORD"
          value_from {
            secret_key_ref {
              name = module.secrets.reddit_password_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name  = "REDDIT_USER_AGENT"
          value = "varity.app@${var.release}"
        }

        env {
          name = "KAFKA_BOOTSTRAP_SERVERS"
          value_from {
            secret_key_ref {
              name = var.kafka_url_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_KEY"
          value_from {
            secret_key_ref {
              name = var.kafka_key_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_SECRET"
          value_from {
            secret_key_ref {
              name = var.kafka_secret_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_SCRAPING_ADDRESS"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.scraping_redis_address.secret
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_SCRAPING_PASSWORD"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.scraping_redis_password.secret
              key  = "latest"
            }
          }
        }
      }
    }
  }

  metadata {
    annotations = {
      "autoscaling.knative.dev/maxScale" = "10"
      generated-by                       = "magic-modules"
      "run.googleapis.com/launch-stage"  = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}

resource "google_cloud_run_service" "scrape_reddit_historical" {
  name     = "scrape-reddit-historical-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 8
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/historical-reddit-scraper:${var.release}"

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }

        env {
          name = "KAFKA_BOOTSTRAP_SERVERS"
          value_from {
            secret_key_ref {
              name = var.kafka_url_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_KEY"
          value_from {
            secret_key_ref {
              name = var.kafka_key_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "KAFKA_AUTH_SECRET"
          value_from {
            secret_key_ref {
              name = var.kafka_secret_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_SCRAPING_ADDRESS"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.scraping_redis_address.secret
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_SCRAPING_PASSWORD"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.scraping_redis_password.secret
              key  = "latest"
            }
          }
        }
      }
    }
  }

  metadata {
    annotations = {
      "autoscaling.knative.dev/maxScale" = "3"
      generated-by                       = "magic-modules"
      "run.googleapis.com/launch-stage"  = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}


resource "google_cloud_run_service" "tiingo" {
  name     = "tiingo-scraper-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 8
      timeout_seconds = 900
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/tiingo:${var.release}"

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }

        env {
          name  = "POSTGRES_ADDRESS"
          value = "/cloudsql/${var.cloud_sql_connection_name}/.s.PGSQL.5432"
        }

        env {
          name  = "POSTGRES_NETWORK"
          value = "unix"
        }

        env {
          name  = "POSTGRES_DB"
          value = "finance"
        }

        env {
          name = "POSTGRES_USERNAME"
          value_from {
            secret_key_ref {
              name = module.secrets.postgres_username_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_PASSWORD"
          value_from {
            secret_key_ref {
              name = module.secrets.postgres_password_secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "TIINGO_TOKEN"
          value_from {
            secret_key_ref {
              name = module.secrets.tiingo_token_secret_id
              key  = "latest"
            }
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "10"
        "run.googleapis.com/cloudsql-instances" = "${var.cloud_sql_connection_name}"
      }
    }
  }

  metadata {
    annotations = {
      generated-by                      = "magic-modules"
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}

resource "google_cloud_run_service" "etl_bigquery_to_influx" {
  name     = "etl-bigquery-to-influx-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 10
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/bigquery-to-influx:${var.release}"

        resources {
          limits = {
            memory = "2Gi"
            cpu = "1000m"
          }
        }

        env {
          name = "LOG_LEVEL"
          value = "DEBUG"
        }

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }

        env {
          name = "INFLUX_URL"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.influx_url.secret
              key  = "latest"
            }
          }
        }

        env {
          name = "INFLUX_TOKEN"
          value_from {
            secret_key_ref {
              name = data.google_secret_manager_secret_version.influx_token.secret
              key  = "latest"
            }
          }
        }

        env {
          name = "INFLUX_BUCKET"
          value = var.deployment
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "10"
      }
    }
  }

  metadata {
    annotations = {
      generated-by                      = "magic-modules"
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}

resource "google_cloud_run_service" "scheduler" {
  name     = "scheduler-${var.deployment}"
  location = var.region
  project  = var.project
  provider = google-beta

  template {
    spec {
      container_concurrency = 50
      containers {
        image = "${var.container_registry}/${var.project}/${var.deployment}/scraping/scheduler:${var.release}"

        resources {
          limits = {
            memory = "2Gi"
            cpu = "1000m"
          }
        }

        env {
          name = "LOG_LEVEL"
          value = "DEBUG"
        }

        env { 
          name = "URL_BIGQUERY_TO_INFLUX"
          value = google_cloud_run_service.etl_bigquery_to_influx.status[0].url
        }

        env {
          name = "DEPLOYMENT_MODE"
          value = var.deployment
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "5"
        "run.googleapis.com/cloudsql-instances" = var.cloud_sql_connection_name
      }
    }
  }

  metadata {
    annotations = {
      generated-by                      = "magic-modules"
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}