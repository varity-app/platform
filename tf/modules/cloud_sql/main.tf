
module "secrets" {
  source = "../secrets"
}

resource "google_sql_database_instance" "main_instance" {
  name    = "varity-sql-${var.deployment}"
  project = var.project
  region  = var.region

  database_version = "POSTGRES_13"

  settings {
    tier      = "db-f1-micro"
    disk_size = 25
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