variable "project_id" {
  description = "The project ID to host the cluster in"
}

variable "region" {
  default = "europe-west4"
}

variable "prefix" {
  description = "Prefix for all resources"
}

variable "main_ip_cidr" {
  default = "172.17.0.0/16"
}

variable "pods_ip_cidr" {
  default = "10.100.0.0/22"
}

variable "svc_ip_cidr" {
  default = "10.1.1.0/24"
}

output "svc_subnet_name" {
  value = "${var.prefix}-svc"
}

output "pods_subnet_name" {
  value = "${var.prefix}-pods"
}

output "network_id" {
  value = google_compute_network.custom.id
}

output "subnetwork_id" {
  value = google_compute_subnetwork.custom.id
}

output "svc_ip_cidr" {
  value = var.svc_ip_cidr
}

output "pods_ip_cidr" {
  value = var.pods_ip_cidr
}

