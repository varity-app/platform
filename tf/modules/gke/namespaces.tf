resource "kubernetes_namespace" "scrapers" {
  metadata {
    name = "scrapers"
  }
}

resource "kubernetes_namespace" "beam" {
  metadata {
    name = "beam"
  }
}