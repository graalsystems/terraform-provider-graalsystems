---
page_title: "GraalSystems: graalsystems_group"
description: |-
Manages GraalSystems Groups.
---

# graalsystems_group

Creates and manages GraalSystems Groups.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Examples

### Basic

```hcl
resource "graalsystems_group" "my_group" {
  name         = "my group"
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the group.

- `description` (Optional) The description of the group.