---
page_title: "Provider: GraalSystems"
description: |-
The GraalSystems provider is used to manage [GraalSystems](https://graal.systems) resources. The provider needs to be configured with the proper credentials before it can be used.
---

# GraalSystems Provider

The GraalSystems provider is used to manage GraalSystems resources.
The provider needs to be configured with the proper credentials before it can be used.

**This is the documentation for the version `>= 1.0.4` of the provider.**

Use the navigation to the left to read about the available resources.

## Terraform 0.13 and later

For Terraform 0.13 and later, please also include this:

```hcl
terraform {
  required_providers {
    graalsystems = {
      source = "graalsystems/graalsystems"
    }
  }
  required_version = ">= 1.0.1"
}
```

## Example

Here is an example that will set up a project with a job, an identity and a library.

You can test this config by creating a `test.tf` and run terraform commands from this directory:

- Get your credentials
- Initialize a Terraform working directory: `terraform init`
- Generate and show the execution plan: `terraform plan`
- Build the infrastructure: `terraform apply`

```hcl
terraform {
  required_providers {
    graalsystems = {
      source = "graalsystems/graalsystems"
      version = "1.0.4"
    }
  }
}

provider "graalsystems" {
  api_url = "https://api.graal.systems"
  auth_url = "https://identity.graal.systems"
  tenant = "XXX"
  username = "XXX"
  password = "XXX"
}

resource "graalsystems_project" "my_project" {
  name = "Example project"
  description = "This is an example project"
}

resource "graalsystems_identity" "my_identity" {
  name = "my identity"
}

resource "graalsystems_job" "my_job" {
  project_id = graalsystems_project.my_project.id
  identity_id = graalsystems_identity.my_identity.id
  name = "my job"
  description = "a useful description job"
  labels = {
    "pipeline": "data-ingestion"
  }

  spark {
    instance_type = "Standard_General_G1_v1"
    main_class_name = "org.acme.MySparkJob"
  }
}
```

## Authentication

The [GraalSystems](https://graal.systems) authentication is based on personal credentials or application credentials.

The GraalSystems provider offers three ways of providing these credentials.
The following methods are supported, in this priority order:

1. [Environment variables](#environment-variables)
1. [Static credentials](#static-credentials)

### Environment variables

!> **Warning**: Not released yet

You can provide your credentials via the `GS_USERNAME`, `GS_PASSWORD` environment variables.

Example:

```hcl
provider "graalsystems" {}
```

Usage:

```bash
$ export GS_USERNAME="my-username"
$ export GS_PASSWORD="my-password"
$ terraform plan
```

### Static credentials

!> **Warning**: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided by adding `access_key` and `secret_key` attributes in-line in the GraalSystems provider block:

Example:

```hcl
provider "graalsystems" {
  username = "XXX"
  password = "XXX"
}
```

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`), the following arguments are supported in the GraalSystems provider block:

| Provider Argument | [Environment Variables](#environment-variables) | Description                                                                                                                             | Mandatory |
|-------------------|-------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|-----------|
| `username`      | `GS_USERNAME`                                | [GraalSystems username](https://console.graal.systems)                                                                 | ✅        |
| `password`      | `GS_PASSWORD`                                | [GraalSystems password](https://console.graal.systems)                                                                 | ✅        |
| `tenant`      | `GS_TENANT`                        | The [tenant ID](https://console.graal.systems/profile) that will be used as default value for all resources.                   | ✅        |
| `api_url`      | `GS_API_URL`                        |                    |         |
| `auth_url`      | `GS_AUTH_URL`                        |     |         |

## Debugging a deployment

In case you want to [debug a deployment](https://www.terraform.io/internals/debugging), you can use the following command to increase the level of verbosity.

`GS_DEBUG=true TF_LOG=WARN TF_LOG_PROVIDER=DEBUG terraform apply`

- `GS_DEBUG`: set the debug level of the graalsystems SDK.
- `TF_LOG`: set the level of the Terraform logging.
- `TF_LOG_PROVIDER`: set the level of the GraalSystems Terraform provider logging.

### Submitting a bug report or a feature request

In case you find something wrong with the graalsystems provider, please submit a bug report on the [Terraform provider repository](https://github.com/graalsystems/terraform-provider-graalsystems/issues/new/choose).
If it is a bug report, please include a **minimal** snippet of the Terraform configuration that triggered the error.
This helps a lot to debug the issue.
