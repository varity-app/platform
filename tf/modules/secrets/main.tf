data "google_secret_manager_secret_version" "reddit_client_id" {
  secret = "varity-reddit-client-id"
}

data "google_secret_manager_secret_version" "reddit_client_secret" {
  secret = "varity-reddit-client-secret"
}

data "google_secret_manager_secret_version" "reddit_username" {
  secret = "varity-reddit-username"
}

data "google_secret_manager_secret_version" "reddit_password" {
  secret = "varity-reddit-password"
}

data "google_secret_manager_secret_version" "reddit_user_agent" {
  secret = "varity-reddit-user-agent"
}

