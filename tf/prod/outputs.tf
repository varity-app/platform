output "metabase_svc_key" {
  value     = module.biquery.metabase_svc_key
  sensitive = true
}
