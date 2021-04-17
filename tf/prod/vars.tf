variable "deployment" {
  description = "General deployment type"
  type        = string
  default     = "prod"
}

variable "release" {
  description = "Varity release"
  type = string
  default = "0.6.0"
}
