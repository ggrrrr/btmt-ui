
module "vpc" {
  source     = "./vpc"
  project_id = var.project_id
  prefix     = var.prefix
}

module "sql" {
  source = "./sql"

  prefix     = var.prefix
  project_id = var.project_id
  region     = var.region

  username      = var.sql_username
  user_password = random_password.sql_password.result

  dev_networkds = var.dev_networkds
  database_name = var.sql_database_name
}

module "gke" {
  source = "./gke"

  prefix     = var.prefix
  region     = var.region
  project_id = var.project_id

  main_network_id   = module.vpc.network_id
  subnetwork_id     = module.vpc.subnetwork_id
  pods_network_name = module.vpc.pods_subnet_name
  svc_network_name  = module.vpc.svc_subnet_name
}


module "k8s" {
  source = "./k8s"

  prefix       = var.prefix
  region       = var.region
  host         = "${var.dns_main}.${var.dns_zone_name}"
  sql_username = var.sql_username
  sql_password = random_password.sql_password.result
  sql_host     = module.sql.public_ip_address
  sql_database = var.sql_database_name
  mgo_uri      = var.mgo_uri
  mgo_username = var.mgo_username
  mgo_password = var.mgo_password
}

output "public_ip" {
  value = module.vpc.public_ip
}

# output "ingress" {
# value = module.k8s.nginx-ingress[0].load_balancer[0].ingress[0].ip
# }
output "ingress-ip" {
  value = module.k8s.nginx-ingress-ip
}
