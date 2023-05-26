package dexidp

import "github.com/hashicorp/terraform-plugin-framework/types"

type providerConfiguration struct {
	Host types.String      `tfsdk:"host"`
	TLS  *tlsConfiguration `tfsdk:"tls"`
}

type tlsConfiguration struct {
	ServerCrt types.String `tfsdk:"ca_crt"`
	ClientCrt types.String `tfsdk:"client_crt"`
	ClientKey types.String `tfsdk:"client_key"`
}
