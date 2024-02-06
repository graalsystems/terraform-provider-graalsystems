---
layout: "graalsystems"
page_title: "GraalSystems: graalsystems_workflow"
description: |-
  Gets information about an existing Workflow.
---

# graalsystems_workflow

Gets information about an existing workflow.

## Example Usage

```hcl
# Get info by ID
data "graalsystems_workflow" "by_id" {
  workflow_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

```hcl
# Get info by name
data "graalsystems_workflow" "by_name" {
  name = "my workflow"
}
```

## Argument Reference

- `name` - (Optional) The name of the workflow.
  Only one of the `name` and `workflow_id` should be specified.

- `workflow_id` - (Optional) The ID of the workflow.
  Only one of the `name` and `workflow_id` should be specified.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

- `description` - The description of the workflow.
- `id` - The ID of the workflow, similar to the `workflow_id` argument.
- `identity_id` - The ID of the identity used to run the workflow.
- `job` - The list of job definitions the workflow chains.
- `labels` - The tag labels of the workflow.
- `name` - The name of the workflow
- `project_id` - The ID of the project where the workflow belongs.
- `schedule` - The workflow schedule definition.
