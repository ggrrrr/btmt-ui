
resource "google_sql_database_instance" "main" {
  name                = "${var.prefix}-sql"
  database_version    = var.database_version
  region              = var.region
  deletion_protection = false

  settings {

    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = var.tier

    ip_configuration {
      ipv4_enabled = true
      # authorized_networks {
      #   name  = "me"
      #   value = "95.42.23.198/32"
      # }
      dynamic "authorized_networks" {
        for_each = var.dev_networkds
        iterator = value
        content {
          name  = "rule ${value.value}"
          value = value.value
        }
      }
    }
  }
}

resource "google_sql_database" "database" {
  name     = var.database_name
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "user" {
  name     = var.username
  password = var.user_password

  instance = google_sql_database_instance.main.name
}


