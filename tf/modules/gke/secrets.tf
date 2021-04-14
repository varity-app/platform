resource "kubernetes_secret" "reddit_credentials" {
    metadata {
        name = "reddit-credentials"
        namespace = kubernetes_namespace.scrapers.metadata[0].name
    }

    data = {
        username = module.secrets.reddit_username
        password = module.secrets.reddit_password
        client_id = module.secrets.reddit_client_id
        client_secret = module.secrets.reddit_client_secret
        user_agent = module.secrets.reddit_user_agent
    }
}

# resource "kubernetes_secret" "google_credentials" {
#     metadata {
#         name = "reddit-credentials"
#         namespace = "scrapers"
#     }

#     data = {
#         username = module.secrets.google_secret_manager_secret_version.reddit_username.plaintext
#         password = module.secrets.google_secret_manager_secret_version.reddit_password.plaintext
#         client_id = module.secrets.google_secret_manager_secret_version.reddit_client_id.plaintext
#         client_secret = module.secrets.google_secret_manager_secret_version.reddit_client_secret.plaintext
#         user_agent = module.secrets.google_secret_manager_secret_version.reddit_user_agent.plaintext
#     }
# }