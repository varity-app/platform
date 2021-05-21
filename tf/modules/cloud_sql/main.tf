terraform {
  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.13.0-pre1"
    }
  }
}

module "secrets" {
  source = "../secrets"
}

resource "google_sql_database" "finance" {
  name     = "finance"
  instance = google_sql_database_instance.main_instance.name
}

resource "google_sql_database_instance" "main_instance" {
  name    = "varity-sql-${var.deployment}"
  project = var.project
  region  = var.region

  database_version = "POSTGRES_11"

  settings {
    tier      = "db-f1-micro"
    disk_size = 50
  }

  deletion_protection = "true"
}

resource "google_sql_user" "main_user" {
  instance = google_sql_database_instance.main_instance.name
  name     = module.secrets.postgres_username
  password = module.secrets.postgres_password
}
