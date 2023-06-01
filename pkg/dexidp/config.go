package dexidp

import "github.com/hashicorp/terraform-plugin-framework/types"

type providerConfiguration struct {
	Host types.String      `tfsdk:"host"`
	TLS  *tlsConfiguration `tfsdk:"tls"`
}

type tlsConfiguration struct {
	ServerCrt types.String `tfsdk:"ca_cert"`
	ClientCrt types.String `tfsdk:"client_cert"`
	ClientKey types.String `tfsdk:"client_key"`
}
