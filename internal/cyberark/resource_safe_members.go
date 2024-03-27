package provider

import (
	"context"
	"fmt"

	cybrapi "github.com/aharriscybr/cybr-api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	htypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &safeMemberResource{}
	_ resource.ResourceWithConfigure = &safeMemberResource{}
)

// NewSafeResource is a helper function to simplify the provider implementation.
func NewSafeMemberResource() resource.Resource {
	return &safeMemberResource{}
}

// msAccountResource is the resource implementation.
type safeMemberResource struct {
	client *cybrapi.Client
}

// Metadata returns the resource type name.
func (r *safeMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_safemembers"
}

type safeMemberModel struct {
	RetentionDays	htypes.Int64 `tfsdk:"retention"`
	RetentionVersions 	htypes.Int64 `tfsdk:"retention_versions"`
	PurgeEnabled	htypes.Bool `tfsdk:"purge"`
	CPM	htypes.String `tfsdk:"cpm_name"`
	Name	htypes.String `tfsdk:"safe_name"`
	Description	htypes.String `tfsdk:"safe_desc"`
	Location	htypes.String `tfsdk:"safe_loc"`
	ID htypes.String `tfsdk:"id"`
	IDNUM htypes.Int64 `tfsdk:"id_number"`
	LastUpdated htypes.String `tfsdk:"last_updated"`
	SeedMember htypes.String `tfsdk:"member"`
	SeedMType htypes.String `tfsdk:"member_type"`
	PermType htypes.String `tfsdk:"permission_level"`
}

func (r *safeMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CyberArk Privilege Cloud Safe Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "CyberArk Privilege Cloud Safe URL ID- Generated from CyberArk after onboarding safe.",
				Computed: true,
			},
			"id_number": schema.Int64Attribute{
				Description: "CyberArk Privilege Cloud Safe ID- Generated from CyberArk after onboarding safe.",
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"safe_name": schema.StringAttribute{
				Description: "The unique name of the Safe. The following characters cannot be used in the Safe name: \\ / : * < > . | ? “% & +",
				Required: true,
			},
			"member": schema.StringAttribute{
				Description: "Owning Safe Member.",
				Required: true,
			},
			"member_type": schema.StringAttribute{
				Description: "Member user type: user or group.",
				Required: true,
			},
			"permission_level": schema.StringAttribute{
				Description: "Membership Permission Level. Currently supported inputs: full, read, approver, manager.",
				Required: true,
			},
			"safe_desc": schema.StringAttribute{
				Description: "The description of the Safe.",
				Optional: true,
			},
			"safe_loc": schema.StringAttribute{
				Description: "The location of the Safe in the Vault.",
				Optional: true,
			},
			"cpm_name": schema.StringAttribute{
				Description: "The name of the CPM user who will manage the new Safe.",
				Optional: true,
			},
			"retention": schema.Int64Attribute{
				Description: "The number of retained versions of every password that is stored in the Safe.",
				Optional: true,
			},
			"retention_versions": schema.Int64Attribute{
				Description: "The number of days that password versions are saved in the Safe.",
				Optional: true,
			},
			"purge": schema.BoolAttribute{
				Description: "Whether or not to automatically purge files after the end of the Object History Retention Period defined in the Safe properties.",
				Optional: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *safeMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create a new resource.
func (r *safeMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	

	
}

// Refresh Existing State
func (r *safeMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	

}


// Update updates the resource and sets the updated Terraform state on success.
func (r *safeMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Update is not supported through terraform. Please consult with your CyberArk Administrator to process account property updates.")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *safeMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete is not supported through terraform. Please consult with your CyberArk Administrator to process deleting this resource.")
}