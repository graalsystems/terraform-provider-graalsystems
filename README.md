# Terraform Provider for GraalSystems

- [Provider Documentation Website](https://www.terraform.io/docs/providers/graalsystems/index.html)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- [![Go Report Card](https://goreportcard.com/badge/github.com/terraform-providers/terraform-provider-graalsystems)](https://goreportcard.com/report/github.com/terraform-providers/terraform-provider-graalsystems)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/graalsystems/terraform-provider-graalsystems`

```sh
$ mkdir -p $GOPATH/src/github.com/graalsystems; cd $GOPATH/src/github.com/graalsystems
$ git clone git@github.com:graalsystems/terraform-provider-graalsystems.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/graalsystems/terraform-provider-graalsystems
$ make build
```

## Using the provider

See the [GraalSystems Provider Documentation](https://registry.terraform.io/providers/graalsystems/graalsystems/latest/docs) to get started using the GraalSystems provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

You have the option to [override](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers) the intended version

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-graalsystems
...
```

Please refer to the [TESTING.md](TESTING.md) for testing.
