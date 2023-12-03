
variable "project_id" {
  description = "The project ID to host the cluster in"
  default     = "bikeareaui"
}

variable "region" {
  default = "europe-west4"
}

variable "prefix" {
  description = "Prefix for all resources"
  default     = "btmt-dev"
}

variable "dev_networkds" {
  type = list(string)
}

variable "sql_username" {
  default = "auth"
}

variable "sql_database_name" {
  default = "auth"
}


provider "google" {
  project = var.project_id
  region  = var.region
}

resource "random_password" "sql_password" {
  length           = 16
  special          = true
  override_special = "_%@"
}
