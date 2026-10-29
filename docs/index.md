---
page_title: "Provider: graal"
description: |-
  The graalsystems provider manages resources of a graal Data Platform installation: projects, identities, jobs, workflows, workspaces, users and groups.
---

# graal Provider

The `graalsystems` provider manages resources of a [graal Data Platform](https://graal.systems/en/docs)
installation through its REST API (`/api/v1`): projects, identities, jobs, workflows,
workspaces, users and groups.

graal is self-hosted: point the provider at your own installation (API URL, OpenID Connect
server and tenant).

Use the navigation to the left to read about the available resources and data sources.

## Example

```hcl
terraform {
  required_providers {
    graalsystems = {
      source  = "graalsystems/graalsystems"
      version = "~> 1.2.1"
    }
  }
}

# Credentials come from GRAAL_CLIENT_ID and GRAAL_CLIENT_SECRET.
provider "graalsystems" {
  api_url  = "https://graal.example.com/api/v1"
  auth_url = "https://sso.example.com"
  tenant   = "acme"
}

resource "graalsystems_project" "analytics" {
  name        = "Analytics"
  description = "Sales analytics"
}

resource "graalsystems_identity" "runner" {
  name = "analytics-runner"
}

resource "graalsystems_job" "nightly_load" {
  name        = "nightly-load"
  project_id  = graalsystems_project.analytics.id
  identity_id = graalsystems_identity.runner.id

  options {
    type          = "python"
    file          = "loader.py"
    instance_type = "Standard_Development_D1_v1"
  }

  library {
    type = "file"
    key  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # graalctl libraries upload loader.py
  }

  schedule {
    type              = "cron"
    cron_expression   = "0 6 * * *"
    timezone          = "Europe/Paris"
    infrastructure_id = "my-infrastructure"
  }

  timeout_seconds = 3600
  max_retries     = 2
}
```

Run `terraform init`, `terraform plan` and `terraform apply` from the directory holding this
configuration.

## Authentication

The provider gets an OpenID Connect token from the installation's Keycloak, in one of two modes:

- `application` (recommended for automation): the id and secret of a graal Application,
  exchanged with the client credentials grant. The application acts with the rights of its
  service account.
- `credentials`: a username and password (resource owner password grant, local accounts only).

The mode is inferred from the attributes that are set; `auth_mode` forces it.

!> **Warning** Do not write secrets in Terraform files. Use environment variables or a
secret manager.

### Environment variables

```bash
export GRAAL_API_URL="https://graal.example.com/api/v1"
export GRAAL_KEYCLOAK_URL="https://sso.example.com"
export GRAAL_TENANT="acme"
export GRAAL_CLIENT_ID="terraform"
export GRAAL_CLIENT_SECRET="..."
terraform plan
```

```hcl
provider "graalsystems" {}
```

The `GS_*` variables of earlier versions (`GS_API_URL`, `GS_AUTH_URL`, `GS_TENANT`,
`GS_APPLICATION_ID`, `GS_APPLICATION_SECRET`, `GS_USERNAME`, `GS_PASSWORD`, `GS_AUTH_MODE`)
are still read.

## Argument Reference

| Argument | Environment variable | Description |
|---|---|---|
| `api_url` | `GRAAL_API_URL` | URL of the REST API, including `/api/v1`. Required. |
| `auth_url` | `GRAAL_KEYCLOAK_URL` | Base URL of the OpenID Connect server (Keycloak), without `/realms/...`. Required. |
| `tenant` | `GRAAL_TENANT` | The tenant, sent as the `X-Tenant` header. Required. |
| `realm` | `GRAAL_REALM` | Authentication realm. Discovered from the tenant when omitted. |
| `auth_mode` | `GRAAL_AUTH_MODE` | `application` or `credentials`. Inferred when omitted. |
| `application_id` | `GRAAL_CLIENT_ID` | Application id (`application` mode). |
| `application_secret` | `GRAAL_CLIENT_SECRET` | Application secret (`application` mode). Sensitive. |
| `username` | `GRAAL_USERNAME` | Username (`credentials` mode). |
| `password` | `GRAAL_PASSWORD` | Password (`credentials` mode). Sensitive. |

## Debugging

```bash
GRAAL_DEBUG=1 TF_LOG=WARN TF_LOG_PROVIDER=DEBUG terraform apply
```

- `GRAAL_DEBUG`: logs the HTTP requests of the graal SDK.
- `TF_LOG`, `TF_LOG_PROVIDER`: [Terraform logging](https://developer.hashicorp.com/terraform/internals/debugging).

## Reporting issues

Open an issue on the [provider repository](https://github.com/graalsystems/terraform-provider-graalsystems/issues)
with a **minimal** configuration that reproduces the problem.
