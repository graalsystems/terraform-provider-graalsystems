---
layout: "graalsystels"
page_title: "GraalSystems: graalsystems_project"
description: |-
Gets information about an existing project.
---

# graalsystems_project

Gets information about an existing project.

## Example Usage

```hcl
# Get info by ID
data graalsystems_project "by_id" {
  project_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Project.
  Only one of the `name` and `project_id` should be specified.

- `project_id` - (Optional) The ID of the Project.
  Only one of the `name` and `project_id` should be specified.
  
