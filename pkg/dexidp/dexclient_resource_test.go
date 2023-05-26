package dexidp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testResourceName = "dexidp_client.test_client"
)

func TestClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "dexidp_client" "test_client" {
	client_id     = "test-client"
	name          = "My Test Client"
	secret        = "Secret"
	redirect_uris = ["https://my-test-app.marcofranssen.nl/callback"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "id", "test-client"),
					resource.TestCheckResourceAttr(testResourceName, "client_id", "test-client"),
					resource.TestCheckResourceAttr(testResourceName, "name", "My Test Client"),
					resource.TestCheckResourceAttr(testResourceName, "secret", "Secret"),
					resource.TestCheckResourceAttr(testResourceName, "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "redirect_uris.0", "https://my-test-app.marcofranssen.nl/callback"),
					resource.TestCheckResourceAttrSet(testResourceName, "last_updated"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "dexidp_client" "test_client" {
	client_id     = "test-client"
	name          = "My Test Client"
	secret        = "Secret"
	redirect_uris = ["https://my-test-app.marcofranssen.nl/callback", "https://another-app.marcofranssen.nl/oidc/callback"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "id", "test-client"),
					resource.TestCheckResourceAttr(testResourceName, "client_id", "test-client"),
					resource.TestCheckResourceAttr(testResourceName, "name", "My Test Client"),
					resource.TestCheckResourceAttr(testResourceName, "secret", "Secret"),
					resource.TestCheckResourceAttr(testResourceName, "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "redirect_uris.0", "https://my-test-app.marcofranssen.nl/callback"),
					resource.TestCheckResourceAttr(testResourceName, "redirect_uris.1", "https://another-app.marcofranssen.nl/oidc/callback"),
					resource.TestCheckResourceAttrSet(testResourceName, "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
