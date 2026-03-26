output "cloud_run_url" {
  value = google_cloud_run_v2_service.main.uri
}

output "db_private_ip" {
  value     = google_sql_database_instance.main.private_ip_address
  sensitive = true
}

output "db_connection_name" {
  value = google_sql_database_instance.main.connection_name
}

output "artifact_registry_url" {
  value = "${var.region}-docker.pkg.dev/${var.project_id}/explainable-engine"
}

output "vpc_connector_name" {
  value = google_vpc_access_connector.connector.name
}
