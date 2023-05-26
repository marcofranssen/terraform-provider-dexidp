package dexidp_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the DexIDP client is properly configured.
	// It is also possible to use the DEXIDP_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "dexidp" {
  host     = "localhost:5557"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"dexidp": providerserver.NewProtocol6WithError(dexidp.New()),
	}
)
