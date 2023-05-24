package dexidp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &dexProvider{}
)

func New() provider.Provider {
	return &dexProvider{}
}

type dexProvider struct{}

// Schema defines the provider-level schema for configuration data.
func (p *dexProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Metadata returns the provider type name.
func (p *dexProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dexidp"
}

// Configure prepares a dexapi gRPC client for data sources and resources.
func (p *dexProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *dexProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *dexProvider) Resources(context.Context) []func() resource.Resource {
	panic("unimplemented")
}
