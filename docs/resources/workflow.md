---
page_title: "GraalSystems: graalsystems_workflow"
description: |-
Manages GraalSystems Workflows.
---

# graalsystems_workflow

Creates and manages GraalSystems Workflows.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Example usage

### Basic

```hcl
data "graalsystems_project" "my_project" {
  project_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}

data "graalsystems_identity" "my_identity" {
  identity_id = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
}

data "graalsystems_job" "my_first_job" {
  job_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}

data "graalsystems_job" "my_second_job" {
  job_id = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}

resource "graalsystems_workflow" "my_workflow" {
  name        = "my workflow"
  description = "my workflow description"
  project_id  = data.graalsystems_project.my_project.id
  identity_id = data.graalsystems_identity.my_identity.id


  job {
    ref  = data.graalsystems_job.my_first_job.id
    name = "First job"
    depends_on = []
  }

  job {
    ref        = data.graalsystems_job.my_second_job.id
    name       = "Second job"
    depends_on = ["First job"]
  }

  schedule {
    type = "once"
  }
}
```

### Advanced

```hcl
data "graalsystems_project" "my_project" {
  project_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}

data "graalsystems_identity" "my_identity" {
  identity_id = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
}

data "graalsystems_job" "my_first_job" {
  job_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}

data "graalsystems_job" "my_second_job" {
  job_id = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}

resource "graalsystems_workflow" "my_workflow" {
  name        = "my workflow"
  description = "my workflow description"
  project_id  = data.graalsystems_project.my_project.id
  identity_id = data.graalsystems_identity.my_identity.id


  job {
    ref        = data.graalsystems_job.my_first_job.id
    name       = "First job"
    depends_on = []
  }

  job {
    ref        = data.graalsystems_job.my_second_job.id
    name       = "Second job"
    depends_on = ["First job"]
  }

  schedule {
    type              = "cron"
    cron_expression   = "0 0 1 1 *"
    timezone          = "Europe/Paris"
    infrastructure_id = "infra-id"
  }

  labels = {
    project = "my project"
  }
}
```

## Arguments Reference

The following arguments are supported:

- `description` (Optional) The description of the workspace.
- `identity_id` (Required) The ID of the identity to use to run the workflow.
- `labels` (Optional) The tag labels of the job.
- `name` - (Required) The name of the workflow.
- `project_id` (Required) The ID of the project to which the workflow belongs.

### job

The job block configures the type of workflow tasks to chain. The definition order is the chaining order.

- `depends_on` (Optional) List of job names (the ones defined in the `name` field) the current job must wait before running.
- `name` (Required) The job name in the workflow.
- `ref` (Required) The job ID to reference.

### schedule

The schedule block configures the schedule of the job. Only one of `cron` or `once` type can be specified.

- `cron_expression` - (Optional) The cron expression to use for the workflow. Only required for `cron` type.
- `device_id` - (Optional) The ID of the device to use for the cron.
- `infrastructure_id` - (Optional) The ID of the infrastructure to use for the workflow. Only required for `cron` type.
- `timezone` - (Optional) The timezone to use for the workflow. Only required for `cron` type.
- `type` - (Required) The type of the schedule.

## Attributes Reference

This resource exports the following attributes in addition to the arguments above:

- `id` - The ID of the workflow.