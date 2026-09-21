package dexidp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testDiscoveryName = "data.dexidp_discovery.test_discovery"
	issuerUrl         = "https://my-issuer.org"
)

func TestDiscoveryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
data "dexidp_discovery" "test_discovery" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDiscoveryName, "issuer", issuerUrl),
					resource.TestCheckResourceAttr(testDiscoveryName, "authorization_endpoint", issuerUrl+"/auth"),
					resource.TestCheckResourceAttr(testDiscoveryName, "token_endpoint", fmt.Sprintf("%s/token", issuerUrl)),
					resource.TestCheckResourceAttr(testDiscoveryName, "jwks_uri", fmt.Sprintf("%s/keys", issuerUrl)),
					resource.TestCheckResourceAttr(testDiscoveryName, "userinfo_endpoint", fmt.Sprintf("%s/userinfo", issuerUrl)),
					resource.TestCheckResourceAttr(testDiscoveryName, "device_authorization_endpoint", fmt.Sprintf("%s/device/code", issuerUrl)),
					resource.TestCheckResourceAttr(testDiscoveryName, "introspection_endpoint", fmt.Sprintf("%s/token/introspect", issuerUrl)),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "grant_types_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "response_types_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "subject_types_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "id_token_signing_alg_values_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "code_challenge_methods_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "scopes_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "token_endpoint_auth_methods_supported.0"),
					resource.TestCheckResourceAttrSet(testDiscoveryName, "claims_supported.0"),
				),
			},
		},
	})
}
