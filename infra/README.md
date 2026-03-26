# Infrastructure — Explainable Engine

## GCP Project

| Property | Value |
|----------|-------|
| Project ID | `explainable-engine-prod` |
| Billing Account | `01A5A4-62E2DB-EDD3BB` (b1) |
| Region | `europe-west1` |
| Artifact Registry | `europe-west1-docker.pkg.dev/explainable-engine-prod/explainable-engine` |

## Service Accounts

| Name | Purpose | Roles |
|------|---------|-------|
| `cloud-run-sa` | Cloud Run runtime | Cloud SQL Client, Secret Manager Accessor |
| `cloud-build-sa` | CI/CD pipeline | Run Developer, Artifact Registry Writer, SA User |

## Enabled APIs

- Cloud Run
- Cloud SQL Admin
- Artifact Registry
- Cloud Build
- Secret Manager
- Service Networking
- Compute Engine
- VPC Access
- Cloud Resource Manager

## Terraform

State stored in GCS bucket: `explainable-engine-tf-state`

```bash
cd infra/terraform
terraform init
terraform plan
terraform apply
```
