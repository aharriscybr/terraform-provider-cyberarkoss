package provider

import (
	"context"
	"fmt"
	"log"
	"time"

	cybrapi "github.com/aharriscybr/cybr-api"
	cybrtypes "github.com/aharriscybr/cybr-api/pkg/cybr/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	htypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &safeObjectResource{}
	_ resource.ResourceWithConfigure = &safeObjectResource{}
)

// NewSafeResource is a helper function to simplify the provider implementation.
func NewSafeResource() resource.Resource {
	return &safeObjectResource{}
}

// msAccountResource is the resource implementation.
type safeObjectResource struct {
	client *cybrapi.Client
}

// Metadata returns the resource type name.
func (r *safeObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_safeobject"
}

type safeObjectModel struct {
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

func (r *safeObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *safeObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *safeObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan safeObjectModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Strings
	var cpm_name, safe_name, safe_desc, safe_loc, member, member_type, permission_level string

	// Boolean
	var purge bool

	// Int -> using int64 to keep consistent with terraform provider type
	var retention, retention_versions int64

	// Not processing env variable overrides for optional fields, this must be defined in the terraform plan.
	// This is a design decision.

	if !plan.Name.IsNull() {
		safe_name = plan.Name.ValueString()
	}

	if !plan.SeedMember.IsNull() {
		member = plan.SeedMember.ValueString()
	}

	if !plan.SeedMType.IsNull() {
		member_type = plan.SeedMType.ValueString()
	}

	if !plan.PermType.IsNull() {
		if ( plan.PermType.ValueString() == "full" || plan.PermType.ValueString() == "read" || plan.PermType.ValueString() == "approver" || plan.PermType.ValueString() == "manager" ) { 
			permission_level = plan.PermType.ValueString()
		} else {
			tflog.Error(ctx, "Permission level does not match acceptable values.")
		}
		
	}

	// Required attributes met
	newSafe := cybrtypes.SafeData {
		Name: &safe_name,
		Owner: &member,
		OwnerType: &member_type,
		Level: &permission_level,
	}

	// Processing optionals
	if !plan.Description.IsNull() {
		safe_desc = plan.Description.ValueString()
		newSafe.Description = &safe_desc
	}

	if !plan.Location.IsNull() {
		safe_loc = plan.Location.ValueString()
		newSafe.Location = &safe_loc
	}

	if !plan.CPM.IsNull() {
		cpm_name = plan.CPM.ValueString()
		newSafe.CPM = &cpm_name
	}

	if !plan.PurgeEnabled.IsNull() {
		purge = plan.PurgeEnabled.ValueBool()
		newSafe.PurgeEnabled = &purge
	}

	if !plan.RetentionDays.IsNull() {
		retention = plan.RetentionDays.ValueInt64()
		newSafe.RetentionDays = &retention
	}

	if !plan.RetentionVersions.IsNull() {
		retention_versions = plan.RetentionVersions.ValueInt64()
		newSafe.RetentionVersions = &retention_versions
	}



	create, err := cybrapi.CreateSafe(&newSafe, r.client.AuthToken, r.client.Domain)
	if err != nil {
		log.Fatal(err)
	}

	if create == nil {

		tflog.Error(ctx, "Error catching ID, this means provisioning failed in pipeline. Please check debug logs and try again.")

	} else {

		plan.ID = htypes.StringValue(*create.URLID)
		plan.IDNUM = htypes.Int64Value(*create.NUMBER)
		plan.LastUpdated = htypes.StringValue(time.Now().Format(time.RFC850))
	
		// Set state to fully populated data
		resp.State.Set(ctx, plan)

	}

	
}

// Refresh Existing State
func (r *safeObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var currState safeObjectModel
	diags := req.State.Get(ctx, &currState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := cybrapi.GetSafe(currState.ID.ValueStringPointer(), r.client.AuthToken, r.client.Domain)
	if err != nil {
		tflog.Error(ctx, "Unable to retrieve state from CyberArk API")
	}

	tflog.Info(ctx, "Refreshing state")

	// Main Refresh Body
	currState.Name = htypes.StringValue(*newState.Name)
	currState.Description = htypes.StringValue(*newState.Description)
	currState.CPM = htypes.StringValue(*newState.CPM)
	currState.Location = htypes.StringValue(*newState.Location)
	currState.RetentionDays = htypes.Int64Value(*newState.RetentionDays)
	currState.PurgeEnabled = htypes.BoolValue(*newState.PurgeEnabled)
	
	// // Set last updated time to last refreshed time
	currState.LastUpdated = htypes.StringValue(time.Now().Format(time.RFC850))

	// Ensure ID is consistent
	currState.ID = htypes.StringValue(*newState.URLID)
	currState.IDNUM = htypes.Int64Value(*newState.NUMBER)

}


// Update updates the resource and sets the updated Terraform state on success.
func (r *safeObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Error(ctx, "Update is not supported through terraform. Please consult with your CyberArk Administrator to process account property updates.")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *safeObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete is not supported through terraform. Please consult with your CyberArk Administrator to process deleting this resource.")
}