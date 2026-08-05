package dexidp_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testResourceName       = "dexidp_client.test_client"
	testPublicResourceName = "dexidp_client.test_public_client"
)

func TestClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
			// ImportState testing
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Dex gRPC
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
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

func TestPublicClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_public_client" {
	client_id     = "test-public-client"
	name          = "My Public Test Client"
	public        = true
	redirect_uris = ["http://localhost:9876/callback"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testPublicResourceName, "id", "test-public-client"),
					resource.TestCheckResourceAttr(testPublicResourceName, "client_id", "test-public-client"),
					resource.TestCheckResourceAttr(testPublicResourceName, "name", "My Public Test Client"),
					resource.TestCheckResourceAttr(testPublicResourceName, "public", "true"),
					resource.TestCheckNoResourceAttr(testPublicResourceName, "secret"),
					resource.TestCheckResourceAttr(testPublicResourceName, "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(testPublicResourceName, "redirect_uris.0", "http://localhost:9876/callback"),
				),
			},
		},
	})
}

func TestClientSecretValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_public_client" {
	client_id     = "test-public-client-empty-secret"
	name          = "My Public Test Client"
	public        = true
	secret        = ""
	redirect_uris = ["http://localhost:9876/callback"]
}
`,
				ExpectError: regexp.MustCompile("Secret Not Allowed"),
			},
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_public_client" {
	client_id     = "test-public-client-secret"
	name          = "My Public Test Client"
	public        = true
	secret        = "must-not-be-configured"
	redirect_uris = ["http://localhost:9876/callback"]
}
`,
				ExpectError: regexp.MustCompile("Secret Not Allowed"),
			},
			{
				Config: GetProviderConfig() + `
resource "dexidp_client" "test_private_client" {
	client_id     = "test-private-client-no-secret"
	name          = "My Private Test Client"
	redirect_uris = ["http://localhost:9876/callback"]
}
`,
				ExpectError: regexp.MustCompile("Secret Required"),
			},
		},
	})
}
