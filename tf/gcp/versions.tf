terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.7.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
