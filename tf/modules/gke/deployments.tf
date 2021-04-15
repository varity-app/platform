resource "kubernetes_deployment" "submissions_scraper" {
  metadata {
    name      = "submissions-scraper"
    namespace = kubernetes_namespace.scrapers.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "submissions-scraper"
      }
    }

    template {

      metadata {
        labels = {
          app = "submissions-scraper"
        }
      }

      spec {
        container {
          image = "cgundlach13/reddit-scraper:${var.release}"
          name  = "scraper"

          resources {
            limits = {
              cpu    = "0.25"
              memory = "256Mi"
            }
            requests = {
              cpu    = "0.125"
              memory = "128Mi"
            }
          }

          env {
            name  = "MODE"
            value = "submissions"
          }

          env {
            name  = "SLEEP"
            value = 120
          }

          env {
            name  = "SUBREDDITS"
            value = var.subreddits
          }

          env {
            name  = "DEPLOYMENT"
            value = var.deployment
          }

          env {
            name = "REDDIT_USERNAME"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "username"
              }
            }
          }

          env {
            name = "REDDIT_PASSWORD"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "password"
              }
            }
          }

          env {
            name = "REDDIT_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "client_id"
              }
            }
          }

          env {
            name = "REDDIT_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "client_secret"
              }
            }
          }

          env {
            name = "REDDIT_USER_AGENT"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "user_agent"
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "comments_scraper" {
  metadata {
    name      = "comments-scraper"
    namespace = kubernetes_namespace.scrapers.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "comments-scraper"
      }
    }

    template {

      metadata {
        labels = {
          app = "comments-scraper"
        }
      }

      spec {
        container {
          image = "cgundlach13/reddit-scraper:${var.release}"
          name  = "scraper"

          resources {
            limits = {
              cpu    = "0.25"
              memory = "256Mi"
            }
            requests = {
              cpu    = "0.125"
              memory = "128Mi"
            }
          }

          env {
            name  = "MODE"
            value = "comments"
          }

          env {
            name  = "SLEEP"
            value = 60
          }

          env {
            name  = "SUBREDDITS"
            value = var.subreddits
          }

          env {
            name  = "DEPLOYMENT"
            value = var.deployment
          }

          env {
            name = "REDDIT_USERNAME"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "username"
              }
            }
          }

          env {
            name = "REDDIT_PASSWORD"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "password"
              }
            }
          }

          env {
            name = "REDDIT_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "client_id"
              }
            }
          }

          env {
            name = "REDDIT_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "client_secret"
              }
            }
          }

          env {
            name = "REDDIT_USER_AGENT"
            value_from {
              secret_key_ref {
                name = "reddit-credentials"
                key  = "user_agent"
              }
            }
          }
        }
      }
    }
  }
}
