---
layout: "graalsystels"
page_title: "GraalSystems: graalsystems_firewall_rule"
description: |-
Gets information about an existing firewall_rule.
---

# graalsystems_firewall_rule

Gets information about an existing firewall_rule.

## Example Usage

```hcl
# Get info by ID
data graalsystems_firewall_rule "by_id" {
  firewall_rule_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Firewall Rule.
  Only one of the `name` and `firewall_rule_id` should be specified.

- `firewall_rule_id` - (Optional) The ID of the Firewall Rule.
  Only one of the `name` and `firewall_rule_id` should be specified.
  
