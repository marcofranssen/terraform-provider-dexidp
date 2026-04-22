package dexidp

import (
	"context"
	"strings"

	"github.com/dexidp/dex/api/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ datasource.DataSource = &clientDataSource{}
)

func NewClientDataSource() datasource.DataSource {
	return &clientDataSource{}
}

func isUnimplementedClientError(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return strings.Contains(err.Error(), "Unimplemented")
	}
	return st.Code() == codes.Unimplemented
}

type clientDataSource struct {
	client api.DexClient
}

type clientDataSourceModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	Name         types.String `tfsdk:"name"`
	LogoURL      types.String `tfsdk:"logo_url"`
	RedirectURIs types.List   `tfsdk:"redirect_uris"`
	TrustedPeers types.List   `tfsdk:"trusted_peers"`
}

func (d *clientDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(api.DexClient)
}

func (d *clientDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (d *clientDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a Dex oauth2 client by ID.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The ID of the Dex oauth2 client.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Dex oauth2 client.",
				Computed:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: "The URL to the logo of the Dex oauth2 client.",
				Computed:    true,
			},
			"redirect_uris": schema.ListAttribute{
				Description: "The allowed redirect_uris for this Dex Oauth2 client.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"trusted_peers": schema.ListAttribute{
				Description: "The trusted peers for this Dex Oauth2 client.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *clientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config clientDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getReq := api.GetClientReq{
		Id: config.ClientID.ValueString(),
	}
	response, err := d.client.GetClient(ctx, &getReq)
	if err != nil {
		if isUnimplementedClientError(err) {
			resp.Diagnostics.AddError(
				"Error getting Dex client",
				"The Dex server does not support the GetClient method. "+
					"This usually means you need to upgrade your Dex server to a newer version. "+
					"The GetClient method was added in Dex API v2 (Dex v2.37+).",
			)
		} else {
			resp.Diagnostics.AddError(
				"Error getting Dex client",
				"Could not get Dex client: "+err.Error(),
			)
		}
		return
	}

	c := response.Client
	config.Name = types.StringValue(c.Name)
	config.LogoURL = types.StringValue(c.LogoUrl)

	redirectURIs, _ := types.ListValueFrom(ctx, types.StringType, c.RedirectUris)
	trustedPeers, _ := types.ListValueFrom(ctx, types.StringType, c.TrustedPeers)
	config.RedirectURIs = redirectURIs
	config.TrustedPeers = trustedPeers

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
