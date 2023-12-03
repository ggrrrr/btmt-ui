
variable "project_id" {
  description = "The project ID to host the cluster in"
}

variable "region" {
  description = "Cloud region"
}

variable "prefix" {
  description = "Prefix for all resources"
}

variable "dev_networkds" {
  type = list(string)
}

variable "database_version" {
  default = "POSTGRES_15"
}

variable "tier" {
  default = "db-f1-micro"
}

variable "database_name" {
  description = "create database with name"
  default     = "auth"
}

variable "username" {
  sensitive = true
}

variable "user_password" {
  sensitive = true
}

output "public_ip_address" {
  value = google_sql_database_instance.main.public_ip_address
}
