output "metabase_svc_key" {
  value     = base64decode(google_service_account_key.metabase_key.private_key)
  sensitive = true
}
