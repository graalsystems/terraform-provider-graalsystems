---
page_title: "GraalSystems: graalsystems_identity"
description: |-
Manages GraalSystems Identities.
---

# graalsystems_identity

Creates and manages GraalSystems Identities.
For more information see [the documentation](https://docs.dev.graal.systems/).

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