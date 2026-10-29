---
page_title: "graal: graalsystems_user"
description: |-
Manages graal Users.
---

# graalsystems_user

Creates and manages graal Users: a local account of the tenant, in the platform and in its
identity server. For more information see [the documentation](https://graal.systems/en/docs).

## Examples

### Basic

```hcl
resource "graalsystems_user" "alice" {
  name       = "alice"
  email      = "alice@example.com"
  first_name = "Alice"
  last_name  = "Martin"
}
```

### With an initial password

```hcl
resource "graalsystems_user" "bob" {
  name     = "bob"
  email    = "bob@example.com"
  password = var.bob_initial_password
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The username of the user.
- `email` - (Required) The email address of the user.
- `first_name` - (Optional) The first name of the user.
- `last_name` - (Optional) The last name of the user.
- `password` - (Optional, Sensitive) Initial password, set at creation only. Without it, the user
  signs in through an identity provider or after a password reset. Changing it afterwards is
  refused: reset the password in the identity server, or recreate the user.
- `description` - (Optional, Deprecated) Ignored: the platform does not store a description for
  users.
