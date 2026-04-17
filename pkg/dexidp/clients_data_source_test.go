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
data "dexidp_clients" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceNameClients, "clients.#"),
				),
			},
		},
	})
}
