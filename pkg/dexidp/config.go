package dexidp

import "github.com/hashicorp/terraform-plugin-framework/types"

type providerConfig struct {
	Host types.String `tfsdk:"host"`
}
