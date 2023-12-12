
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
variable "prefix" {
}
variable "region" {
}
variable "host" {
}
variable "sql_host" {
}
variable "sql_username" {
}
variable "sql_password" {
}
variable "sql_database" {
}
variable "mgo_uri" {
}
variable "mgo_username" {
}
variable "mgo_password" {
}
variable "dns_contact_email" {

}
