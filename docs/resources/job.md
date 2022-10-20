---
page_title: "GraalSystems: graalsystems_job"
description: |-
Manages GraalSystems Jobs.
---

# graalsystems_job

Creates and manages GraalSystems Jobs.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Examples

### Basic

```hcl
resource "graalsystems_job" "my_job" {
  name         = "my job"
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the job.

- `description` (Optional) The description of the job.