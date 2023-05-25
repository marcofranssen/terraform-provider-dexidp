package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp"
)

// Provider documentation generation.
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name dexidp
func main() {
	providerserver.Serve(context.Background(), dexidp.New, providerserver.ServeOpts{
		Address: "hashicorp.com/marcofranssen/dexidp",
	})
}
