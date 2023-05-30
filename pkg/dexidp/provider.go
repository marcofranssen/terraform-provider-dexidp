package dexidp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/dexidp/dex/api/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp/client"
)

var (
	_ provider.Provider = &dexProvider{}
)

func New() provider.Provider {
	return &dexProvider{}
}

type dexProvider struct{}

// Metadata returns the provider type name.
func (p *dexProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dexidp"
}

// Schema defines the provider-level schema for configuration data.
func (p *dexProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provision Dex resources via the Dex gRPC API.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The host and port of the Dex gRPC API.",
				Required:    true,
			},
			"tls": schema.SingleNestedAttribute{
				Description: "TLS configuration for the Dex gRPC API.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"ca_crt": schema.StringAttribute{
						Description: "The server certificate authority for the Dex gRPC API",
						Optional:    true,
					},
					"client_crt": schema.StringAttribute{
						Description: "Client certificate for mutual TLS authentication",
						Optional:    true,
					},
					"client_key": schema.StringAttribute{
						Description: "Client key for mutual TLS authentication",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}
}

// Configure prepares a dexapi gRPC client for data sources and resources.
func (p *dexProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerConfiguration
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Dex IDP API host",
			"The provider cannot create the Dex IDP API client as there is an unknown configuration value for the Dex IDP API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DEXIDP_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("DEXIDP_HOST")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Dex IDP API host",
			"The provider cannot create the Dex IDP API client as there is a missing or empty value for the Dex IDP API host. "+
				"Set the host value in the configuration or use the DEXIDP_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := constructDexClient(host, config.TLS)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Dex IDP API Client",
			"An unexpected error occurred when creating the Dex IDP API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Dex IDP Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func constructDexClient(host string, tlsCfg *tlsConfiguration) (api.DexClient, error) {
	if tlsCfg == nil {
		return client.New(host, insecure.NewCredentials())
	}

	caCrt := tlsCfg.ServerCrt.ValueString()
	clientCrt := tlsCfg.ClientCrt.ValueString()
	clientKey := tlsCfg.ClientKey.ValueString()

	if (clientCrt == "" && clientKey != "") || (clientCrt != "" && clientKey == "") {
		return nil, fmt.Errorf("client crt and key must both be set or both be empty")
	}

	cPool := x509.NewCertPool()

	if !cPool.AppendCertsFromPEM([]byte(caCrt)) {
		return nil, fmt.Errorf("failed to parse CA crt")
	}

	clientCert, err := tls.X509KeyPair([]byte(clientCrt), []byte(clientKey))
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key pair: %v", err)
	}

	clientTLSConfig := &tls.Config{
		RootCAs:      cPool,
		Certificates: []tls.Certificate{clientCert},
	}

	return client.New(host, credentials.NewTLS(clientTLSConfig))
}

// DataSources defines the data sources implemented in the provider.
func (p *dexProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *dexProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDexClientResource,
	}
}
