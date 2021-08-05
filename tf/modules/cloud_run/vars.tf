variable "project" {
  description = "GCP project name"
  type        = string
  default     = "varity"
}

variable "deployment" {
  description = "General deployment type"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-east1"
}

variable "cloud_sql_connection_name" {
  description = "Connection name belonging to the Cloud SQL instance"
  type        = string
}

variable "kafka_url_secret_id" {
  description = "Kafka bootstrap servers URL secret id"
  type        = string
}

variable "kafka_key_secret_id" {
  description = "SASL credentials for authenticating to kafka secret id"
  type        = string
}

variable "kafka_secret_secret_id" {
  description = "SASL credentials for authenticating to kafka secret id"
  type        = string
}