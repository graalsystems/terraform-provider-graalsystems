---
layout: "graalsystels"
page_title: "GraalSystems: graalsystems_identity"
description: |-
Gets information about an existing identity.
---

# graalsystems_identity

Gets information about an existing identity.

## Example Usage

```hcl
# Get info by ID
data graalsystems_identity "by_id" {
  identity_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Identity.
  Only one of the `name` and `identity_id` should be specified.

- `identity_id` - (Optional) The ID of the Identity.
  Only one of the `name` and `identity_id` should be specified.
  
