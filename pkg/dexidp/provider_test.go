package dexidp_test

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp"
)

func GetProviderConfig() string {
	cwd, _ := os.Getwd()

	return fmt.Sprintf(`
provider "dexidp" {
	host = "127.0.0.1:5557"
	tls = {
		ca_cert     = file("%s/../../certs/ca.crt")
		client_cert = file("%s/../../certs/client.crt")
		client_key  = file("%s/../../certs/client.key")
	}
}`, cwd, cwd, cwd)
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"dexidp": providerserver.NewProtocol6WithError(dexidp.New()),
	}
)
