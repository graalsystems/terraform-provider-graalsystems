---
layout: "graalsystems"
page_title: "GraalSystems: graalsystems_group"
description: |-
Gets information about an existing group.
---

# graalsystems_group

Gets information about an existing group.

## Example Usage

```hcl
# Get info by ID
data graalsystems_group "by_id" {
  group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Group.
  Only one of the `name` and `group_id` should be specified.

- `group_id` - (Optional) The ID of the Group.
  Only one of the `name` and `group_id` should be specified.
  
