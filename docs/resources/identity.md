---
page_title: "graal: graalsystems_identity"
description: |-
Manages graal Identities.
---

# graalsystems_identity

Creates and manages graal Identities.
For more information see [the documentation](https://graal.systems/en/docs).

## Examples

### Basic

```hcl
resource "graalsystems_identity" "my_identity" {
  name         = "my identity"
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the identity.

- `description` (Optional) The description of the identity.