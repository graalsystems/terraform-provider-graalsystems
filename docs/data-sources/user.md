---
layout: "graalsystels"
page_title: "GraalSystems: graalsystems_user"
description: |-
Gets information about an existing user.
---

# graalsystems_user

Gets information about an existing user.

## Example Usage

```hcl
# Get info by ID
data graalsystems_user "by_id" {
  user_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the User.
  Only one of the `name` and `user_id` should be specified.

- `user_id` - (Optional) The ID of the User.
  Only one of the `name` and `user_id` should be specified.
  
