---
page_title: "GraalSystems: graalsystems_user"
description: |-
Manages GraalSystems Users.
---

# graalsystems_user

Creates and manages GraalSystems Users.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Examples

### Basic

```hcl
resource "graalsystems_user" "my_user" {
  name         = "my user"
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the user.

- `description` (Optional) The description of the user.