
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

output "public_ip" {
  value = module.vpc.public_ip
}
