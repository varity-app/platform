variable "project" {
  description = "GCP project name"
  type        = string
  default     = "varity"
}

variable "deployment" {
  description = "General deployment type"
  type        = string
}

variable "location" {
    description = "location of the GKE cluster"
    type = string
    default = "us-east1-c"
}