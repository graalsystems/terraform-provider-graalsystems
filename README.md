# Terraform provider for graal

The `graalsystems` provider manages resources of a [graal Data Platform](https://graal.systems/en/docs)
installation — projects, identities, jobs, workflows, workspaces, users and groups — through its
REST API (`/api/v1`).

- [Provider documentation](docs/index.md) (also on the
  [Terraform Registry](https://registry.terraform.io/providers/graalsystems/graalsystems/latest/docs))
- [Example configuration](examples/main.tf)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) 1.0 or later
  (or [OpenTofu](https://opentofu.org))
- A graal Data Platform installation, and the id and secret of a graal Application with the
  rights needed on the managed projects
- To build the provider: [Go](https://go.dev/doc/install) (version in [go.mod](go.mod))

## Using the provider

```hcl
terraform {
  required_providers {
    graalsystems = {
      source  = "graalsystems/graalsystems"
      version = "~> 1.2.2"
    }
  }
}

provider "graalsystems" {
  api_url  = "https://graal.example.com/api/v1"
  auth_url = "https://sso.example.com"
  tenant   = "acme"
  # application_id / application_secret from GRAAL_CLIENT_ID / GRAAL_CLIENT_SECRET
}
```

See [docs/index.md](docs/index.md) for authentication and all provider arguments.

## Developing the provider

```sh
go build ./...
go test ./...           # unit tests, against an in-process fake of the API
go vet ./...
```

To try a local build, add a [development override](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "graalsystems/graalsystems" = "/path/to/go/bin"
  }
  direct {}
}
```

then run `go install` and use Terraform as usual (skip `terraform init`).

## Releasing

Push a `v*` tag (for example `v1.2.0`), or run the Release workflow by hand from `main` with
the version (Actions → Release → Run workflow), which creates the tag. The release workflow builds the provider for all
platforms with GoReleaser, signs the checksums with the GPG key stored in the `GPG_PRIVATE_KEY`
and `PASSPHRASE` secrets, and creates a draft GitHub release. Publish it (Releases page, or the
Publish release workflow) to make the version available on the Terraform Registry.

## License

Apache License 2.0, see [LICENSE](LICENSE).
