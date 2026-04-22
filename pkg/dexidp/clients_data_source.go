package dexidp

import (
	"context"

	"github.com/dexidp/dex/api/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &clientsDataSource{}
)

func NewClientsDataSource() datasource.DataSource {
	return &clientsDataSource{}
}

type clientsDataSource struct {
	client api.DexClient
}

type clientInfoModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	Name         types.String `tfsdk:"name"`
	LogoURL      types.String `tfsdk:"logo_url"`
	RedirectURIs types.List   `tfsdk:"redirect_uris"`
	TrustedPeers types.List   `tfsdk:"trusted_peers"`
}

type clientsDataSourceModel struct {
	Clients []clientInfoModel `tfsdk:"clients"`
}

func (d *clientsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(api.DexClient)
}

func (d *clientsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clients"
}

func (d *clientsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all Dex oauth2 clients.",
		Attributes: map[string]schema.Attribute{
			"clients": schema.ListNestedAttribute{
				Description: "List of all Dex oauth2 clients.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"client_id": schema.StringAttribute{
							Description: "The ID of the Dex oauth2 client.",
							Computed:    true,
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
				},
				Computed: true,
			},
		},
	}
}

func (d *clientsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config clientsDataSourceModel

	listReq := api.ListClientReq{}
	response, err := d.client.ListClients(ctx, &listReq)
	if err != nil {
		if isUnimplementedClientError(err) {
			resp.Diagnostics.AddError(
				"Error listing Dex clients",
				"The Dex server does not support the ListClients method. "+
					"This usually means you need to upgrade your Dex server to a newer version.",
			)
		} else {
			resp.Diagnostics.AddError(
				"Error listing Dex clients",
				"Could not list Dex clients: "+err.Error(),
			)
		}
		return
	}

	var clients []clientInfoModel
	for _, c := range response.Clients {
		redirectURIs, _ := types.ListValueFrom(ctx, types.StringType, c.RedirectUris)
		trustedPeers, _ := types.ListValueFrom(ctx, types.StringType, c.TrustedPeers)

		clients = append(clients, clientInfoModel{
			ClientID:     types.StringValue(c.Id),
			Name:         types.StringValue(c.Name),
			LogoURL:      types.StringValue(c.LogoUrl),
			RedirectURIs: redirectURIs,
			TrustedPeers: trustedPeers,
		})
	}

	config.Clients = clients
	diags := resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
