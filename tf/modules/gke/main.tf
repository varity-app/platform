# Retrieve an access token as the Terraform runner
data "google_client_config" "provider" {}

resource "google_container_cluster" "varity_cluster" {
  name     = "varity-k8s-${var.deployment}"
  location = var.location
  project  = var.project

  remove_default_node_pool = true
  initial_node_count       = 1
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

