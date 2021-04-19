variable "release" {
  description = "Varity release"
  type        = string
}

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
  type        = string
  default     = "us-east1-c"
}

variable "main_email" {
  description = "Email associated with google account"
  type        = string
  default     = "gundlachcallum@gmail.com"
}

variable "subreddits" {
  description = "Subreddits to scrape"
  type        = string
  default     = "wallstreetbets,smallstreetbets,stocks,valueinvesting,securityanalysis,investing"
}

variable "container_repository" {
  description = "URL and prefix of the container repository"
  type = string
  default = "gcr.io/varity/scraping"
}
