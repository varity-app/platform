data "google_secret_manager_secret_version" "reddit_client_id" {
  secret  = "varity-reddit-client-id"
  project = var.project
}

data "google_secret_manager_secret_version" "reddit_client_secret" {
  secret  = "varity-reddit-client-secret"
  project = var.project
}

data "google_secret_manager_secret_version" "reddit_username" {
  secret  = "varity-reddit-username"
  project = var.project
}

data "google_secret_manager_secret_version" "reddit_password" {
  secret  = "varity-reddit-password"
  project = var.project
}

data "google_secret_manager_secret_version" "reddit_user_agent" {
  secret  = "varity-reddit-user-agent"
  project = var.project
}

data "google_secret_manager_secret_version" "postgres_username" {
  secret  = "varity-postgres-username"
  project = var.project
}

data "google_secret_manager_secret_version" "postgres_password" {
  secret  = "varity-postgres-password"
  project = var.project
}