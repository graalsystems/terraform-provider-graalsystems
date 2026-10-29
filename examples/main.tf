# Example: a project with a scheduled Python job and a two-step workflow.
#
#   export GRAAL_API_URL=https://graal.example.com/api/v1
#   export GRAAL_KEYCLOAK_URL=https://sso.example.com
#   export GRAAL_TENANT=acme
#   export GRAAL_CLIENT_ID=terraform
#   export GRAAL_CLIENT_SECRET=...
#   terraform init && terraform apply

terraform {
  required_providers {
    graalsystems = {
      source  = "graalsystems/graalsystems"
      version = "~> 1.2"
    }
  }
}

provider "graalsystems" {}

variable "infrastructure_id" {
  description = "Infrastructure the scheduled runs are placed on (graalctl infrastructures list)."
  type        = string
}

variable "code_library_key" {
  description = "Key of the uploaded code (graalctl libraries upload extract.py)."
  type        = string
}

resource "graalsystems_project" "analytics" {
  name        = "Analytics"
  description = "Sales analytics"
}

resource "graalsystems_identity" "runner" {
  name = "analytics-runner"
}

resource "graalsystems_job" "extract" {
  name        = "extract"
  project_id  = graalsystems_project.analytics.id
  identity_id = graalsystems_identity.runner.id

  options {
    type          = "python"
    file          = "extract.py"
    instance_type = "Standard_Development_D1_v1"
  }

  library {
    type = "file"
    key  = var.code_library_key
  }

  timeout_seconds = 1800
  max_retries     = 2
}

resource "graalsystems_job" "report" {
  name        = "report"
  project_id  = graalsystems_project.analytics.id
  identity_id = graalsystems_identity.runner.id

  options {
    type          = "bash"
    lines         = ["echo \"Report for $GRAAL_LOGICAL_DATE\""]
    instance_type = "Standard_Development_D0_v1"
  }
}

resource "graalsystems_workflow" "daily" {
  name        = "daily"
  project_id  = graalsystems_project.analytics.id
  identity_id = graalsystems_identity.runner.id

  job {
    ref        = graalsystems_job.extract.id
    name       = "extract"
    depends_on = []
  }

  job {
    ref        = graalsystems_job.report.id
    name       = "report"
    depends_on = ["extract"]
  }

  schedule {
    type              = "cron"
    cron_expression   = "0 6 * * *"
    timezone          = "Europe/Paris"
    infrastructure_id = var.infrastructure_id
  }
}
