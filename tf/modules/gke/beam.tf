resource "kubernetes_deployment" "comments_executor" {
  metadata {
    name      = "comments-executor"
    namespace = kubernetes_namespace.beam.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "comments-executor"
      }
    }

    template {

      metadata {
        labels = {
          app = "comments-executor"
        }
      }

      spec {
        container {
          image             = "${var.container_repository}/${var.deployment}/beam-executor:${var.release}"
          image_pull_policy = "Always"
          name              = "beam"

          resources {
            limits = {
              cpu    = "0.25"
              memory = "512Mi"
            }
            requests = {
              cpu    = "0.125"
              memory = "256Mi"
            }
          }

          env {
            name  = "PIPELINE"
            value = "comments"
          }

          env {
            name  = "DEPLOYMENT"
            value = var.deployment
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "submissions_executor" {
  metadata {
    name      = "submissions-executor"
    namespace = kubernetes_namespace.beam.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "submissions-executor"
      }
    }

    template {

      metadata {
        labels = {
          app = "submissions-executor"
        }
      }

      spec {
        container {
          image             = "${var.container_repository}/${var.deployment}/beam-executor:${var.release}"
          image_pull_policy = "Always"
          name              = "beam"

          resources {
            limits = {
              cpu    = "0.25"
              memory = "512Mi"
            }
            requests = {
              cpu    = "0.125"
              memory = "256Mi"
            }
          }

          env {
            name  = "PIPELINE"
            value = "submissions"
          }

          env {
            name  = "DEPLOYMENT"
            value = var.deployment
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "scraped_posts_executor" {
  metadata {
    name      = "scraped-posts-executor"
    namespace = kubernetes_namespace.beam.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "scraped-posts-executor"
      }
    }

    template {

      metadata {
        labels = {
          app = "scraped-posts-executor"
        }
      }

      spec {
        container {
          image             = "${var.container_repository}/${var.deployment}/beam-executor:${var.release}"
          image_pull_policy = "Always"
          name              = "beam"

          resources {
            limits = {
              cpu    = "0.25"
              memory = "512Mi"
            }
            requests = {
              cpu    = "0.125"
              memory = "256Mi"
            }
          }

          env {
            name  = "PIPELINE"
            value = "scraped_posts"
          }

          env {
            name  = "DEPLOYMENT"
            value = var.deployment
          }
        }
      }
    }
  }
}
