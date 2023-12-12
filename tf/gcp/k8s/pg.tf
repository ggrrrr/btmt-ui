
resource "kubernetes_secret" "example" {
  metadata {
    name = "postgres"
  }

  data = {
    POSTGRES_USERNAME = var.sql_username
    POSTGRES_PASSWORD = var.sql_password
    POSTGRES_HOST     = var.sql_host
    POSTGRES_DATABASE = var.sql_database
  }

}

resource "kubernetes_secret" "mgo" {
  metadata {
    name = "mgo"
  }

  data = {
    MGO_DATABASE = "test"
    MGO_USERNAME = var.mgo_username
    MGO_PASSWORD = var.mgo_password
    MGO_URI      = var.mgo_uri
  }

}
