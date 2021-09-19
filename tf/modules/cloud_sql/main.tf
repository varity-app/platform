
module "secrets" {
  source = "../secrets"
}

resource "google_sql_database_instance" "main_instance" {
  name    = "varity-sql-${var.deployment}"
  project = var.project
  region  = var.region

  database_version = "POSTGRES_13"

  settings {
    tier      = "db-g1-small"
    disk_size = 25

    database_flags {
      name = "max_connections"
      value = 1000
    }
  }

  deletion_protection = "true"
}

# This postgres DB stores data from IEX Cloud
resource "google_sql_database" "finance" {
  name     = "finance"
  instance = google_sql_database_instance.main_instance.name
}

# Admin user
resource "google_sql_user" "main_user" {
  instance = google_sql_database_instance.main_instance.name
  name     = module.secrets.postgres_username
  password = module.secrets.postgres_password
}