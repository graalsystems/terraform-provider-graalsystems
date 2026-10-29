---
page_title: "graal: graalsystems_workspace"
description: |-
Manages graal Workspaces.
---

# graalsystems_workspace

Creates and manages graal Workspaces.
For more information see [the documentation](https://graal.systems/en/docs).

## Example usage

### Basic

```hcl
resource "graalsystems_workspace" "my_workspace" {
  name        = "my workspace"
  description = "my description"
  type        = "vscode"
  
  infrastructure_id = "infra-id"
  instance_type     = "Standard_Development_D1_v1"
}
```

## Arguments Reference

The following arguments are supported:

- `description` (Optional) The description of the workspace.
- `infrastructure_id` - (Required) The ID of the infrastructure the workspace will be deployed to.
- `instance_type` - (Required) The compute instance type used to run the workspace.
- `name` - (Required) The name of the workspace.
- `type` (Required) The type of workspace: `jupyter`, `vscode` (code-server), `rstudio`, `grafana`, `superset`, `metabase` or `zeppelin`.

## Attributes Reference

This resource exports the following attributes in addition to the arguments above:

- `id` - The ID of the workspace.
- `owner` - The owner ID of the workspace.
- `status` - The status of the workspace.
- `version` - The version of the workspace according to its type.
- `public_url` - The URL to access the workspace.

Changing `type`, `infrastructure_id` or `instance_type` deletes the workspace and creates a new one.
