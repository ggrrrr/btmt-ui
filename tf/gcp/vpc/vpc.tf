resource "google_compute_network" "custom" {
  name                    = var.prefix
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "custom" {
  name          = var.prefix
  ip_cidr_range = var.main_ip_cidr
  region        = var.region
  network       = google_compute_network.custom.id
  secondary_ip_range {
    range_name    = "${var.prefix}-svc"
    ip_cidr_range = var.svc_ip_cidr
  }

  secondary_ip_range {
    range_name    = "${var.prefix}-pods"
    ip_cidr_range = var.pods_ip_cidr
  }
}
