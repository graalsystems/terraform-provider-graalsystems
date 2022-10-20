---
page_title: "GraalSystems: graalsystems_project"
description: |-
Manages GraalSystems Projects.
---

# graalsystems_project

Creates and manages GraalSystems Projects.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Examples

### Basic

```hcl
resource "graalsystems_project" "my_project" {
  name         = "my project"
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the project.

- `description` (Optional) The description of the project.