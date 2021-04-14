resource "kubernetes_namespace" "scrapers" {
    metadata {
        name = "scrapers"
    }
}