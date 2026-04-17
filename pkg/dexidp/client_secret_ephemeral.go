package dexidp

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource = &clientSecretEphemeral{}
)

func NewClientSecretEphemeral() ephemeral.EphemeralResource {
	return &clientSecretEphemeral{}
}

type clientSecretEphemeral struct{}

type clientSecretEphemeralModel struct {
	ClientID types.String `tfsdk:"client_id"`
	Secret   types.String `tfsdk:"secret"`
	Result   types.String `tfsdk:"result"`
}

func (e *clientSecretEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_secret"
}

func (e *clientSecretEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Ephemeral client secret for passing to resources without persisting to state.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The ID of the Dex oauth2 client.",
				Optional:    true,
			},
			"secret": schema.StringAttribute{
				Description: "The secret for the Dex oauth2 client.",
				Optional:    true,
				Sensitive:   true,
			},
			"result": schema.StringAttribute{
				Description: "The ephemeral result containing the secret.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (e *clientSecretEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var model clientSecretEphemeralModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret := model.Secret.ValueString()
	if secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			"Missing secret",
			"Either secret must be provided or the ephemeral resource must be used with an ephemeral input.",
		)
		return
	}

	model.Result = types.StringValue(secret)

	diags = resp.Result.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)

	resp.RenewAt = time.Now().Add(24 * time.Hour)
}
