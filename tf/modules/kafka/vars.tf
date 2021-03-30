variable "bootstrap_servers" {
  description = "URL of the Confluent Cloud boostrap servers"
  type        = string
}

variable "confluent_key" {
  description = "SASL username to login to the confluent cluster"
  type        = string
}

variable "confluent_secret" {
  description = "SASL password to login to the confluent cluster"
  type        = string
}

variable "num_partitions" {
  description = "Number of partitions to use for each topic"
  type        = number
}
