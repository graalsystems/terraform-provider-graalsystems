package main

import (
	"github.com/graalsystems/terraform-provider-graalsystems/graalsystems"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: graalsystems.Provider(graalsystems.DefaultProviderConfig()),
	})

}
