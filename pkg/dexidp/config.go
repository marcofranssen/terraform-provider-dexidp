package dexidp

import "github.com/hashicorp/terraform-plugin-framework/types"

type providerConfiguration struct {
	Host types.String `tfsdk:"host"`
}
