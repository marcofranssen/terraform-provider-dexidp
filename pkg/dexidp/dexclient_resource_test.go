package dexidp_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testResourceName         = "dexidp_client.test_client"
	testResourceNamePublic   = "dexidp_client.test_client_public"
	testResourceNameNoSecret = "dexidp_client.test_client_no_secret"
)

func TestClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create non-public client (should work with secret)
			{
				Config: GetProviderConfig() + `
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
			// Create public client without secret (should work)
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_client_public" {
	client_id     = "test-client-public"
	name          = "My Public Client"
	public        = true
	redirect_uris = ["https://my-public-app.marcofranssen.nl/callback"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceNamePublic, "id", "test-client-public"),
					resource.TestCheckResourceAttr(testResourceNamePublic, "client_id", "test-client-public"),
					resource.TestCheckResourceAttr(testResourceNamePublic, "name", "My Public Client"),
					resource.TestCheckResourceAttr(testResourceNamePublic, "public", "true"),
					resource.TestCheckResourceAttr(testResourceNamePublic, "redirect_uris.#", "1"),
					resource.TestCheckNoResourceAttr(testResourceNamePublic, "secret"),
				),
			},
			// Create non-public client without secret (should fail with validation error)
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_client_no_secret" {
	client_id     = "test-client-no-secret"
	name          = "My No Secret Client"
	public       = false
	redirect_uris = ["https://my-no-secret-app.marcofranssen.nl/callback"]
}
`,
				ExpectError: regexp.MustCompile("secret is required"),
			},
			// ImportState testing
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}
