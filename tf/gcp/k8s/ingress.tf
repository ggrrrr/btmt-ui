
provider "helm" {
  kubernetes {
    host  = "https://${data.google_container_cluster.gke.endpoint}"
    token = data.google_client_config.provider.access_token
    cluster_ca_certificate = base64decode(
      data.google_container_cluster.gke.master_auth[0].cluster_ca_certificate,
    )
  }
}

resource "helm_release" "nginx_ingress" {
  name             = "nginx-ingress-controller"
  namespace        = "nginx-ingress"
  create_namespace = true

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"

  #   set {
  #     name  = "namespace"
  #     value = "nginx-ingress"
  #   }
}

resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true

  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"

  # version    = "v1.18"
  set {
    name  = "installCRDs"
    value = "true"
  }
}


data "kubernetes_service_v1" "nginx-ingress" {
  metadata {
    namespace = "nginx-ingress"
    name      = "nginx-ingress-controller"
  }
}

output "nginx-ingress-ip" {
  #   value = data.kubernetes_service_v1.nginx-ingress.status
  value = data.kubernetes_service_v1.nginx-ingress.status[0].load_balancer[0].ingress[0].ip
  #   value = data.kubernetes_service_v1.nginx-ingress.status[0].load_balancer[0].ingress[0].ip
}

resource "kubernetes_manifest" "test-crd" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"

    metadata = {
      name = "letsencrypt-dev"
      #   namespace = "cert-manager"
      labels = {
        name = "letsencrypt-dev"
      }
    }

    spec = {
      acme = {
        email = "ggrrrr@gmail.com"
        privateKeySecretRef = {
          name = "letsencrypt-dev"
        }
        server = "https://acme-v02.api.letsencrypt.org/directory"
        solvers = [
          {
            "http01" = {
              "ingress" = {
                "class" = "nginx"
              }
            }
          }
        ]
      }
    }
  }
}
