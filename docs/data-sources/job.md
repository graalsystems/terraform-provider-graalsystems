---
layout: "graalsystems"
page_title: "GraalSystems: graalsystems_job"
description: |-
Gets information about an existing job.
---

# graalsystems_job

Gets information about an existing job.

## Example Usage

```hcl
# Get info by ID
data graalsystems_job "by_id" {
  job_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Job.
  Only one of the `name` and `job_id` should be specified.

- `job_id` - (Optional) The ID of the Job.
  Only one of the `name` and `job_id` should be specified.
  
