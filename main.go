package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/dennismdejong/terraform-provider-dexidp/pkg/dexidp"
)

func main() {
	providerserver.Serve(context.Background(), dexidp.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/dennismdejong/dexidp",
	})
}