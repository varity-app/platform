resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "primary-preemtible-node-pool"
  project    = var.project
  location   = var.location
  cluster    = google_container_cluster.varity_cluster.name
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-standard-4"

    taint {
        key = "cloud.google.com/gke-preemptible"
        value = true
        effect = "NO_SCHEDULE"
    }

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "primary-node-pool"
  project    = var.project
  location   = var.location
  cluster    = google_container_cluster.varity_cluster.name
  node_count = 1

  node_config {
    # preemptible  = true
    machine_type = "e2-standard-2"

    # Google recommends custom service accounts that have cloud-platform scope and 
    # permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}