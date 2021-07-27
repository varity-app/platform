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