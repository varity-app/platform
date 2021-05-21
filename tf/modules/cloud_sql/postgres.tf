provider "postgresql" {
  scheme   = "gcppostgres"
  host     = google_sql_database_instance.public_ip_address
  database = google_sql_database.finance.name
  username = module.secrets.postgres_username
  password = module.secrets.postgres_password
  //  sslmode         = "require"
}


