---
page_title: "graal: graalsystems_group"
description: |-
Manages graal Groups.
---

# graalsystems_group

Creates and manages graal Groups.
For more information see [the documentation](https://graal.systems/en/docs).

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

The API cannot edit groups: changing `name` or `description` deletes the group and creates a new one.
