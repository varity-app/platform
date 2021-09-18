variable "project" {
  description = "GCP project name"
  type        = string
  default     = "varity"
}

variable "deployment" {
  description = "General deployment type"
  type        = string
}

variable "etl_bigquery_to_influx_url" {
  description = "URL of the bigquery -> influx ETL service runnning on Cloud Run"
  type        = string
}
