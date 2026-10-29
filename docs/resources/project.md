---
page_title: "graal: graalsystems_project"
description: |-
Manages graal Projects.
---

# graalsystems_project

Creates and manages graal Projects.
For more information see [the documentation](https://graal.systems/en/docs).

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