---
page_title: "graal: graalsystems_user"
description: |-
Manages graal Users.
---

# graalsystems_user

Creates and manages graal Users.
For more information see [the documentation](https://graal.systems/en/docs).

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

- `name` - (Required) The username of the user.

- `description` (Optional) The description of the user.