
variable "cloudflare_api_token" {

}

variable "dns_zone_name" {

}

variable "dns_main" {
  default = "dev-gcp"
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "dns_zone" {
  name = var.dns_zone_name
}

output "dns_zone" {
  value = data.cloudflare_zone.dns_zone
}

resource "cloudflare_record" "example" {
  zone_id = data.cloudflare_zone.dns_zone.id
  name    = "dev-gcp"
  value   = module.k8s.nginx-ingress-ip
  type    = "A"
  ttl     = 3600
}
