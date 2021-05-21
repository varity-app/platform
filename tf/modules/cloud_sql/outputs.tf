output "postgres_ip" {
  value = google_sql_database_instance.main_instance.public_ip_address
}