output "reddit_client_id" {
  value     = data.google_secret_manager_secret_version.reddit_client_id.secret_data
  sensitive = true
}

output "reddit_client_secret" {
  value     = data.google_secret_manager_secret_version.reddit_client_secret.secret_data
  sensitive = true
}

output "reddit_user_agent" {
  value     = data.google_secret_manager_secret_version.reddit_user_agent.secret_data
  sensitive = true
}

output "reddit_username" {
  value     = data.google_secret_manager_secret_version.reddit_username.secret_data
  sensitive = true
}

output "reddit_password" {
  value     = data.google_secret_manager_secret_version.reddit_password.secret_data
  sensitive = true
}

output "postgres_username" {
  value     = data.google_secret_manager_secret_version.postgres_username.secret_data
  sensitive = true
}

output "postgres_password" {
  value     = data.google_secret_manager_secret_version.postgres_password.secret_data
  sensitive = true
}
