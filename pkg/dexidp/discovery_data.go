package dexidp

import (
	"context"
	"fmt"
	"github.com/dexidp/dex/api/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &dexDiscoveryDataSource{}
)

// NewDexClientResource instantiates a new Dex Client resource.
func NewDexDiscoveryDataSource() datasource.DataSource {
	return &dexDiscoveryDataSource{}
}

type dexDiscoveryDataSource struct {
	client api.DexClient
}

type dexDiscoveryModel struct {
	Issuer                            types.String `tfsdk:"issuer"`
	AuthorizationEndpoint             types.String `tfsdk:"authorization_endpoint"`
	TokenEndpoint                     types.String `tfsdk:"token_endpoint"`
	JwksUri                           types.String `tfsdk:"jwks_uri"`
	UserinfoEndpoint                  types.String `tfsdk:"userinfo_endpoint"`
	DeviceAuthorizationEndpoint       types.String `tfsdk:"device_authorization_endpoint"`
	IntrospectionEndpoint             types.String `tfsdk:"introspection_endpoint"`
	GrantTypesSupported               types.List   `tfsdk:"grant_types_supported"`
	ResponseTypesSupported            types.List   `tfsdk:"response_types_supported"`
	SubjectTypesSupported             types.List   `tfsdk:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported  types.List   `tfsdk:"id_token_signing_alg_values_supported"`
	CodeChallengeMethodsSupported     types.List   `tfsdk:"code_challenge_methods_supported"`
	ScopesSupported                   types.List   `tfsdk:"scopes_supported"`
	TokenEndpointAuthMethodsSupported types.List   `tfsdk:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   types.List   `tfsdk:"claims_supported"`
}

func (r *dexDiscoveryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(api.DexClient)
}

// Metadata returns the resource type name.
func (r *dexDiscoveryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_discovery"
}

// Schema defines the schema for the resource.
func (r *dexDiscoveryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Dex Discovery information.",
		Attributes: map[string]schema.Attribute{
			"issuer": schema.StringAttribute{
				Computed: true,
			},
			"authorization_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"token_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"jwks_uri": schema.StringAttribute{
				Computed: true,
			},
			"userinfo_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"device_authorization_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"introspection_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"grant_types_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"response_types_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"subject_types_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"id_token_signing_alg_values_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"code_challenge_methods_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"scopes_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"token_endpoint_auth_methods_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"claims_supported": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dexDiscoveryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dexDiscoveryModel
	diags := resp.State.Get(ctx, &state)

	disoveryReq := api.DiscoveryReq{}
	response, err := r.client.GetDiscovery(ctx, &disoveryReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Dex discovery",
			fmt.Sprintf("Could not get Dex discovery, unexpected error: %v", err),
		)
		return
	}
	c := response

	state.Issuer = types.StringValue(c.Issuer)
	state.AuthorizationEndpoint = types.StringValue(c.AuthorizationEndpoint)
	state.TokenEndpoint = types.StringValue(c.TokenEndpoint)
	state.JwksUri = types.StringValue(c.JwksUri)
	state.UserinfoEndpoint = types.StringValue(c.UserinfoEndpoint)
	state.DeviceAuthorizationEndpoint = types.StringValue(c.DeviceAuthorizationEndpoint)
	state.IntrospectionEndpoint = types.StringValue(c.IntrospectionEndpoint)
	state.GrantTypesSupported, _ = types.ListValueFrom(ctx, types.StringType, c.GrantTypesSupported)
	state.ResponseTypesSupported, _ = types.ListValueFrom(ctx, types.StringType, c.ResponseTypesSupported)
	state.SubjectTypesSupported, _ = types.ListValueFrom(ctx, types.StringType, c.SubjectTypesSupported)
	state.IdTokenSigningAlgValuesSupported, _ = types.ListValueFrom(ctx, types.StringType, c.IdTokenSigningAlgValuesSupported)
	state.CodeChallengeMethodsSupported, _ = types.ListValueFrom(ctx, types.StringType, c.CodeChallengeMethodsSupported)
	state.ScopesSupported, _ = types.ListValueFrom(ctx, types.StringType, c.ScopesSupported)
	state.TokenEndpointAuthMethodsSupported, _ = types.ListValueFrom(ctx, types.StringType, c.TokenEndpointAuthMethodsSupported)
	state.ClaimsSupported, _ = types.ListValueFrom(ctx, types.StringType, c.ClaimsSupported)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
