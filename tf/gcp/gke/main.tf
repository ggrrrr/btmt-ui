
variable "project_id" {
  description = "The project ID to host the cluster in"
}

variable "region" {
  description = "Cloud region"
}

variable "prefix" {
  description = "Prefix for all resources"
}

variable "main_network_id" {
  description = "VPC network id"
}

variable "subnetwork_id" {
  description = "VPC subnetwork id"
}

variable "pods_network_name" {
  description = "VPC network for pods IPs"
}

variable "svc_network_name" {
  description = "VPC network for services IPs"
}

output "endpoint" {
  value = google_container_cluster.gke.endpoint
}
