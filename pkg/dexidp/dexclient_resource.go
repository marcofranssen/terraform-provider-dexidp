package dexidp

import (
	"context"
	"fmt"
	"time"

	"github.com/dexidp/dex/api/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dexClientResoure{}
	_ resource.ResourceWithConfigure = &dexClientResoure{}
)

// NewDexClientResource instantiates a new Dex Client resource.
func NewDexClientResource() resource.Resource {
	return &dexClientResoure{}
}

type dexClientResoure struct {
	client api.DexClient
}

type dexClientModel struct {
	ID           types.String `tfsdk:"id"`
	ClientID     types.String `tfsdk:"client_id"`
	Secret       types.String `tfsdk:"secret"`
	Name         types.String `tfsdk:"name"`
	Public       types.Bool   `tfsdk:"public"`
	LogoURL      types.String `tfsdk:"logo_url"`
	RedirectURIs types.List   `tfsdk:"redirect_uris"`
	TrustedPeers types.List   `tfsdk:"trusted_peers"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}

// Configure adds the provider configured client to the resource.
func (r *dexClientResoure) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(api.DexClient)
}

// Metadata returns the resource type name.
func (r *dexClientResoure) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

// Schema defines the schema for the resource.
func (r *dexClientResoure) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provision a Dex oauth2 client.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of your terraform resoure (Set to client_id automatically).",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Dex client.",
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The ID of your Dex oauth2 client.",
				Required:    true,
			},
			"secret": schema.StringAttribute{
				Description: "The Secret of your Dex oauth2 client.",
				Required:    true,
				Sensitive:   true,
			},
			"public": schema.BoolAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of your Dex oauth2 client.",
				Required:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: "The url to the logo of your Dex oauth2 client.",
				Optional:    true,
			},
			"redirect_uris": schema.ListAttribute{
				Description: "The allowed redirect_uris for this Dex Oauth2 client.",
				ElementType: types.StringType,
				Required:    true,
			},
			"trusted_peers": schema.ListAttribute{
				Description: "The trusted peers for this Dex Oauth2 client.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dexClientResoure) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dexClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectURIs := utils.ListStringValuesToSlice(plan.RedirectURIs)
	trustedPeers := utils.ListStringValuesToSlice(plan.TrustedPeers)

	createClientReq := api.CreateClientReq{
		Client: &api.Client{
			Id:           plan.ClientID.ValueString(),
			Secret:       plan.Secret.ValueString(),
			Name:         plan.Name.ValueString(),
			Public:       plan.Public.ValueBool(),
			RedirectUris: redirectURIs,
			TrustedPeers: trustedPeers,
			LogoUrl:      plan.LogoURL.ValueString(),
		},
	}
	response, err := r.client.CreateClient(ctx, &createClientReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Dex client",
			fmt.Sprintf("Could not create Dex client, unexpected error: %v", err),
		)
		return
	}
	if response.AlreadyExists {
		resp.Diagnostics.AddError(
			"Error creating Dex client",
			fmt.Sprintf("Could not create Dex client, client with this Name already exists: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(response.Client.GetId())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dexClientResoure) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dexClientResoure) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dexClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectURIs := utils.ListStringValuesToSlice(plan.RedirectURIs)
	trustedPeers := utils.ListStringValuesToSlice(plan.TrustedPeers)

	updateClientReq := api.UpdateClientReq{
		Id:           plan.ClientID.ValueString(),
		Name:         plan.Name.ValueString(),
		RedirectUris: redirectURIs,
		TrustedPeers: trustedPeers,
		LogoUrl:      plan.LogoURL.ValueString(),
	}

	_, err := r.client.UpdateClient(ctx, &updateClientReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Dex Client",
			fmt.Sprintf("Could not update client, unexpected error: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(plan.ClientID.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dexClientResoure) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dexClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteClientReq := api.DeleteClientReq{
		Id: state.ClientID.ValueString(),
	}

	_, err := r.client.DeleteClient(ctx, &deleteClientReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Dex Client",
			"Could not delete client, unexpected error: "+err.Error(),
		)
		return
	}
}
