---
page_title: "graal: graalsystems_job"
description: |-
Manages graal Jobs.
---

# graalsystems_job

Creates and manages graal Jobs.
For more information see [the documentation](https://graal.systems/en/docs).

~> **NOTE:** The only supported job types in this version are `bash` and `python`.

## Example usage

### Basic

```hcl
resource "graalsystems_project" "my_project" {
  name        = "my project"
  description = "my description"
}

resource "graalsystems_identity" "my_identity" {
  name = "my_identity"
}

resource "graalsystems_job" "my_job" {
  name        = "my job"
  description = "my description"
  project_id  = graalsystems_project.my_project.id
  identity_id = graalsystems_identity.my_identity.id

  options {
    type          = "bash"
    lines         = ["echo 'Hello World!'"]
    instance_type = "Standard_Development_D0_v1"
  }
}
```

### Advanced

```hcl
resource "graalsystems_project" "my_project" {
  name        = "my project"
  description = "my description"
}

resource "graalsystems_identity" "my_identity" {
  name = "my_identity"
}

resource "graalsystems_job" "my_job" {
  name        = "my job"
  description = "my description"
  project_id  = graalsystems_project.my_project.id
  identity_id = graalsystems_identity.my_identity.id

  options {
    type          = "python"
    module        = "package.module"
    instance_type = "Standard_Development_D0_v1"
  }
  library {
    type = "file"
    key  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
  schedule {
    type              = "cron"
    cron_expression   = "0 0 * * *"
    timezone          = "Europe/Paris"
    infrastructure_id = "infra-id"
  }

  timeout_seconds = 3600
  max_retries     = 2

  labels = {
    project = "my project"
  }
}
```

## Arguments Reference

The following arguments are supported:

- `description` (Optional) The description of the job.
- `max_retries` (Optional) Number of automatic restarts after a failed run (0 or more).
- `parameters` (Optional) List of parameters passed to the job.
- `secrets` (Optional) List of project secret ids delivered to the runs.
- `identity_id` (Required) The ID of the identity to use to run the job.
- `labels` (Optional) The tag labels of the job.
- `library` (Optional) The library configuration to specify the library to use in the job.
- `name` - (Required) The name of the job.
- `options` - (Required) The options configuration indicates the type of job.
- `project_id` - (Required) The ID of the project to which the job belongs.
- `schedule` - (Optional) The schedule configuration of the job. Without it, the job only runs when started from the console, the API, `graalctl` or an agent.
- `timeout_seconds` (Optional) Maximum duration of a run, in seconds.

### options

The options block configures the job type. Depending on the type, different options are available.

- `docker_image` - (Optional) Container image to use. Defaults to the runtime image configured on the installation for the job type (`toolbox` for bash, `python` for python).
- `env` - (Optional) Environment variables of the runs.
- `instance_type` - (Required) The instance type (CPU, memory, GPU template) of the runs. List them with `graalctl instance-types list`.
- `lines` - (Optional) The bash lines to execute. Only required for `bash` type.
- `module` - (Optional) The python module to execute. Only required for `python` type. Equivalent to `python -m <module>`.
- `type` - (Required) The type of the job.

### schedule

The schedule block configures the schedule of the job. Only one of `cron` or `once` type can be specified.

- `cron_expression` - (Optional) The cron expression to use for the job. Only required for `cron` type.
- `device_id` - (Optional) The ID of the device to use for the cron.
- `infrastructure_id` - (Optional) The ID of the infrastructure to use for the job. Only required for `cron` type.
- `timezone` - (Optional) The timezone to use for the job. Only required for `cron` type.
- `type` - (Required) The type of the schedule.

### library

The library block configures the library to use for the job.
You can specify multiple libraries by defining multiple `library` blocks.

- `key` - (Required) The ID of the library to use for the job.
- `type` - (Required) The type of the library.

## Attributes Reference

This resource exports the following attributes in addition to the arguments above:

- `id` - The ID of the job.