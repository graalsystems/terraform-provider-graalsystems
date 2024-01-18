---
page_title: "GraalSystems: graalsystems_workspace"
description: |-
Manages GraalSystems Workspaces.
---

# graalsystems_workspace

Creates and manages GraalSystems Workspaces.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Example usage

### Basic

```hcl
resource "graalsystems_workspace" "my_workspace" {
  name        = "my workspace"
  description = "my description"
  type        = "vscode"
  
  infrastructure_id = "infra-id"
  instance_type     = "t3.medium"
}
```

## Arguments Reference

The following arguments are supported:

- `description` (Optional) The description of the workspace.
- `infrastructure_id` - (Required) The ID of the infrastructure the workspace will be deployed to.
- `instance_type` - (Required) The compute instance type used to run the workspace.
- `name` - (Required) The name of the workspace.
- `type` (Required) The type of workspace to deploy.

## Attributes Reference

This resource exports the following attributes in addition to the arguments above:

- `id` - The ID of the workspace.
- `owner` - The owner ID of the workspace.
- `status` - The status of the workspace.
- `version` - The version of the workspace according to its type.
- `public_url` - The URL to access the workspace.