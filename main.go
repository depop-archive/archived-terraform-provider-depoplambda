package main

import (
	depoplambda "github.com/depop/terraform-provider-depoplambda/depoplambda"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: depoplambda.Provider})
}
