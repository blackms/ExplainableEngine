variable "project_id" {
  description = "GCP project ID"
  type        = string
  default     = "explainable-engine-prod"
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "europe-west1"
}

variable "db_tier" {
  description = "Cloud SQL machine type"
  type        = string
  default     = "db-f1-micro"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "explainable_engine"
}

variable "db_user" {
  description = "Database user"
  type        = string
  default     = "app_user"
}
