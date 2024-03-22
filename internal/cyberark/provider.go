package provider

import (
	"context"
	"os"

	cybrapi "github.com/aharriscybr/cybr-api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	htypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &cyberarkProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cyberarkProvider {
			version: version,
		}
	}
}

type cyberarkProvider struct {
	version string
}

type confModel struct {
	Tenant htypes.String `tfsdk:"tenant"`
	ClientID htypes.String `tfsdk:"clientid"`
	ClientSecret htypes.String `tfsdk:"clientsecret"`
	Domain htypes.String `tfsdk:"domain"`
}

// Metadata returns the provider type name.
func (p *cyberarkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cyberarkoss"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *cyberarkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	
	resp.Schema = schema.Schema{
		Description: "Configure tenant used to onboard account types into CyberArk Privilege Cloud Vault",
		Attributes: map[string]schema.Attribute{
			"tenant": schema.StringAttribute{
				Description: "CyberArk Shared Services Tenant.",
				Required: true,
			},
			"clientid": schema.StringAttribute{
				Description: "CyberArk Client ID, formatted as username@cyberark.cloud.tenant.",
				Required: true,
			},
			"clientsecret": schema.StringAttribute{
				Description: "CyberArk Client ID Password.",
				Required: true,
				Sensitive: true,
			},
			"domain": schema.StringAttribute{
				Description: "CyberArk Privilege Cloud Domain.",
				Required: true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *cyberarkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	
	var hconfig confModel

	// this little magic right here
	// maps our config from tf file to the local hashi model tfsdk which maps to our go type
	// which then translates to OUR data ... interesting
	resp.Diagnostics.Append(req.Config.Get(ctx, &hconfig)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var t, cid, csec, do string
	t = os.Getenv("CYBERARK_PROVIDER_TENANT")
	cid = os.Getenv("CYBERARK_PROVIDER_CLIENT_ID")
	csec = os.Getenv("CYBERARK_PROVIDER_CLIENT_SECRET")
	do = os.Getenv("CYBERARK_PROVIDER_DOMAIN")

	// Override
	if !hconfig.Tenant.IsNull() {
		t = hconfig.Tenant.ValueString()
	}
	if !hconfig.ClientID.IsNull() {
		cid = hconfig.ClientID.ValueString()
	}
	if !hconfig.ClientSecret.IsNull() {
		csec = hconfig.ClientSecret.ValueString()
	}
	if !hconfig.Domain.IsNull() {
		do = hconfig.Domain.ValueString()
	}

	client, err := cybrapi.NewClient(&t, &do, &cid, &csec)
	if err != nil {
		tflog.Error(ctx, "Error configuring new client.")
		return
	}

	tflog.Info(ctx, "Configured client.")
	resp.DataSourceData = client
	resp.ResourceData = client

}

// DataSources defines the data sources implemented in the provider.
func (p *cyberarkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource {
		NewTokenDataSource,
	  }
}

// Resources defines the resources implemented in the provider.
func (p *cyberarkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDBAccountResource,
		NewAWSAccountResource,
		NewMSAccountResource,
		NewSafeResource,
	}
	
}