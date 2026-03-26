resource "google_cloud_run_v2_service" "main" {
  name     = "explainable-engine"
  location = var.region

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 10
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/explainable-engine/server:latest"

      ports {
        container_port = 8000
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      env {
        name  = "STORAGE_BACKEND"
        value = "postgresql"
      }

      env {
        name  = "DB_HOST"
        value = google_sql_database_instance.main.private_ip_address
      }

      env {
        name  = "DB_PORT"
        value = "5432"
      }

      env {
        name  = "DB_NAME"
        value = var.db_name
      }

      env {
        name  = "DB_USER"
        value = var.db_user
      }

      env {
        name = "DB_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "DB_SSLMODE"
        value = "disable"
      }

      env {
        name  = "LOG_LEVEL"
        value = "INFO"
      }

      env {
        name  = "CORS_ORIGINS"
        value = "*"
      }

      startup_probe {
        http_get {
          path = "/health"
          port = 8000
        }
        initial_delay_seconds = 5
        period_seconds        = 3
        failure_threshold     = 10
      }

      liveness_probe {
        http_get {
          path = "/health"
          port = 8000
        }
        period_seconds = 30
      }
    }

    vpc_access {
      connector = google_vpc_access_connector.connector.id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    service_account = "cloud-run-sa@${var.project_id}.iam.gserviceaccount.com"

    max_instance_request_concurrency = 80
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

# Allow unauthenticated access (public API)
resource "google_cloud_run_v2_service_iam_member" "public" {
  name     = google_cloud_run_v2_service.main.name
  location = var.region
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Grant Cloud Run SA access to Secret Manager secret
resource "google_secret_manager_secret_iam_member" "cloud_run_access" {
  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:cloud-run-sa@${var.project_id}.iam.gserviceaccount.com"
}
