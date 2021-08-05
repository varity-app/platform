output "postgres_ip" {
  value = google_sql_database_instance.main_instance.public_ip_address
}

output "connection_name" {
  value = google_sql_database_instance.main_instance.connection_name
}