package dexidp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testDataSourceNameClients = "data.dexidp_clients.test"
)

func TestClientsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test" {
  client_id    = "test-client"
  name         = "Test Client"
  public       = true
  redirect_uris = ["http://localhost/callback"]
}

data "dexidp_clients" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.#", "1"),
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.0.client_id", "test-client"),
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.0.name", "Test Client"),
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.0.public", "true"),
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(testDataSourceNameClients, "clients.0.redirect_uris.0", "http://localhost/callback"),
				),
			},
		},
	})
}
