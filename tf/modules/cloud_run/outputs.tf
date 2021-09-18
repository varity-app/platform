output "etl_bigquery_to_influx_url" {
	value = google_cloud_run_service.etl_bigquery_to_influx.status[0].url
}