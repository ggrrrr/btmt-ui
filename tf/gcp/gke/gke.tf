
// https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster
resource "google_container_cluster" "gke" {
  name     = var.prefix
  location = var.region

  deletion_protection      = false
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = var.main_network_id
  subnetwork = var.subnetwork_id
  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_network_name
    services_secondary_range_name = var.svc_network_name
  }
}

resource "google_container_node_pool" "linux_pool" {
  name     = "${var.prefix}-lnx-pool"
  project  = google_container_cluster.gke.project
  cluster  = google_container_cluster.gke.name
  location = google_container_cluster.gke.location

  autoscaling {
    max_node_count = "3"
    min_node_count = "1"
  }

  node_config {
    machine_type = "e2-medium"
    image_type   = "COS_CONTAINERD"
    disk_size_gb = 30
    # spot         = true
  }

}
