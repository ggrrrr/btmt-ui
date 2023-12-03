
data "google_client_config" "provider" {}

data "google_container_cluster" "gke" {
  name       = var.prefix
  location   = var.region
  depends_on = [module.gke]
}

provider "kubernetes" {
  # host  = "https://${google_container_cluster.gke.endpoint}"
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
    username = var.sql_username
    password = random_password.sql_password.result
    host     = module.sql.public_ip_address
    database = var.sql_database_name
  }
}


locals {
  object_json = "{\"BASE_URL\":\"http://${module.vpc.public_ip}/rest\"}"
}

output "asd" {
  value = local.object_json
}

resource "kubernetes_config_map" "web-config-json" {
  metadata {
    name = "web-config-json"
  }
  data = {
    "config.json" = local.object_json
  }
}
