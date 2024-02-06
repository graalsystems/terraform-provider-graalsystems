---
layout: "graalsystems"
page_title: "GraalSystems: graalsystems_workspace"
description: |-
Gets information about an existing workspace.
---

# graalsystems_workspace

Gets information about an existing workspace.

## Example Usage

```hcl
# Get info by ID
data graalsystems_workspace "by_id" {
  workspace_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```
    
```hcl
# Get info by name
data graalsystems_workspace "by_name" {
  name = "my workspace"
}
```

## Argument Reference

- `name` - (Optional) The name of the workspace.
  Only one of the `name` and `workpace_id` should be specified.

- `workpace_id` - (Optional) The ID of the workspace.
  Only one of the `name` and `workpace_id` should be specified.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

- `description` - The description of the workspace.
- `id` - The ID of the workspace, similar to the `workspace_id` argument.
- `infrastructure_id` - The ID of the infrastructure the workspace is deployed to.
- `instance_type` - The compute instance type used to run the workspace.
- `name` - The name of the workspace
- `owner` - The owner ID of the workspace.
- `public_url` - The URL to access the workspace.
- `status` - The status of the workspace.
- `type` - The type of workspace.
- `version` - The version of the workspace according to its type.
