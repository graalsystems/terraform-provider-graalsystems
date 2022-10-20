---
page_title: "GraalSystems: graalsystems_firewall_rule"
description: |-
Manages GraalSystems Firewall Rules.
---

# graalsystems_firewall_rule

Creates and manages GraalSystems Firewall Rules.
For more information see [the documentation](https://docs.dev.graal.systems/).

## Examples

### Basic

```hcl
resource "graalsystems_firewall_rule" "my_firewall_rule" {
  description  = "my description"
}

```

## Arguments Reference

The following arguments are supported:

- `description` (Optional) The description of the firewall_rule.