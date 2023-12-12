locals {
  object_json = "{\"BASE_URL\":\"https://${var.host}/rest\"}"
}

resource "kubernetes_config_map" "web-config-json" {
  metadata {
    name = "web-config-json"
  }
  data = {
    "config.json" = local.object_json
  }
}
