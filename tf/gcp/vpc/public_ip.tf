resource "google_compute_global_address" "public_ip" {
  project      = var.project_id
  name         = "${var.prefix}-public-ip"
  address_type = "EXTERNAL"
  ip_version   = "IPV4"
}

output "public_ip" {
  value = google_compute_global_address.public_ip.address
}

output "public_ip_name" {
  value = google_compute_global_address.public_ip.name
}
