resource "google_artifact_registry_repository" "btmt" {
  location      = var.region
  repository_id = "btmt"
  description   = "app docker images"
  format        = "docker"
}

output "registry_host" {
  value = "${google_artifact_registry_repository.btmt.location}-docker.pkg.dev"
}

output "registry_uri" {
  value = "/${google_artifact_registry_repository.btmt.project}/${google_artifact_registry_repository.btmt.name}"
}
