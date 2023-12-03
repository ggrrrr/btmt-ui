

data "google_client_config" "provider" {}

data "google_container_cluster" "gke" {
  name     = var.prefix
  location = var.region
  #   depends_on = [module.gke]
}

provider "kubernetes" {
  host  = "https://${data.google_container_cluster.gke.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.gke.master_auth[0].cluster_ca_certificate,
  )
}

resource "kubernetes_secret" "example" {
  metadata {
    name = "postgres"
  }

  data = {
    username = var.username
    password = var.password
    host     = var.host
    database = var.database_name
  }

}

variable "prefix" {
}
variable "region" {
}
variable "username" {
}

variable "host" {
}

variable "password" {

}
variable "database_name" {

}
