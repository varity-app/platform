# Retrieve an access token as the Terraform runner
data "google_client_config" "provider" {}

resource "google_container_cluster" "varity_cluster" {
  name     = "varity-k8s-${var.deployment}"
  location = var.location
  project  = var.project

  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "primary-node-pool"
  project    = var.project
  location   = var.location
  cluster    = google_container_cluster.varity_cluster.name
  node_count = 1

  node_config {
    # preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and 
    # permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

provider "kubernetes" {
  host  = "https://${google_container_cluster.varity_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    google_container_cluster.varity_cluster.master_auth[0].cluster_ca_certificate,
  )
}

module "secrets" {
  source = "../secrets"
}

