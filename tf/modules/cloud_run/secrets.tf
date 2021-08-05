data "google_project" "project" {
  #   provider = google-beta
}

resource "google_secret_manager_secret_iam_member" "reddit_client_id_access" {
  #   provider = google-beta

  secret_id = module.secrets.reddit_client_id_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "reddit_client_secret_access" {
  #   provider = google-beta

  secret_id = module.secrets.reddit_client_secret_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "reddit_user_agent_access" {
  #   provider = google-beta

  secret_id = module.secrets.reddit_user_agent_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "reddit_username_access" {
  #   provider = google-beta

  secret_id = module.secrets.reddit_username_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "reddit_password_access" {
  #   provider = google-beta

  secret_id = module.secrets.reddit_password_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "postgres_password_access" {
  #   provider = google-beta

  secret_id = module.secrets.postgres_password_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "postgres_username_access" {
  #   provider = google-beta

  secret_id = module.secrets.postgres_username_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "kafka_url_access" {
  #   provider = google-beta

  secret_id = var.kafka_url_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "kafka_key_access" {
  #   provider = google-beta

  secret_id = var.kafka_key_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "kafka_secret_access" {
  #   provider = google-beta

  secret_id = var.kafka_secret_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
}