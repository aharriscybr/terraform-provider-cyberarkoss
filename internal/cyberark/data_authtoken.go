package provider

import (
	"context"
	"fmt"

	// CyberArk Includes
	cybrapi "github.com/aharriscybr/cybr-api"

	// Hashi Includes

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	htypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tokenDataSource{}
	_ datasource.DataSourceWithConfigure = &tokenDataSource{}
)

// NewTokenDataSource is a helper function to simplify the provider implementation.
func NewTokenDataSource() datasource.DataSource {
	return &tokenDataSource{}
}

// dbAccountResource is the resource implementation.
type tokenDataSource struct {
	client *cybrapi.Client
}

// Metadata returns the resource type name.
func (d *tokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authtoken"
}

type tokenDataSourceModel struct {

	Token htypes.String `tfsdk:"token"`

}

func (d *tokenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Shared Services Auth Token",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Shared Services Authorization Token",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to this datasource.
func (d *tokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cybrapi.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cybrapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),

		)

		return
	}

	d.client = client
}

// Return token
func (d *tokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state tokenDataSourceModel

	state.Token = htypes.StringValue(*d.client.AuthToken)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
