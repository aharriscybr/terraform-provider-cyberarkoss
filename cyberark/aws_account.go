package provider

import (
	"context"
	"fmt"
	"log"
	"os"
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
	_ resource.Resource              = &awsAccountResource{}
	_ resource.ResourceWithConfigure = &awsAccountResource{}
)

// NewAWSAccountResource is a helper function to simplify the provider implementation.
func NewAWSAccountResource() resource.Resource {
	return &awsAccountResource{}
}

// awsAccountResource is the resource implementation.
type awsAccountResource struct {
	client *cybrapi.Client
}

// Metadata returns the resource type name.
func (r *awsAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_awsaccount"
}

type awsCredModel struct {
	Name 		htypes.String `tfsdk:"name"`
	Username 	htypes.String `tfsdk:"username"`
	Platform 	htypes.String `tfsdk:"platform"`
	Safe 		htypes.String `tfsdk:"safe"`
	SecretType 	htypes.String `tfsdk:"secrettype"`
	Secret 		htypes.String `tfsdk:"secret"`
	ID 			htypes.String `tfsdk:"id"`
	LastUpdated htypes.String `tfsdk:"last_updated"`
	Manage 		htypes.Bool `tfsdk:"sm_manage"`
	ManageReason	htypes.String `tfsdk:"sm_manage_reason"`
	AWSKID 		htypes.String `tfsdk:"aws_kid"`
	AWSAccount 		htypes.String `tfsdk:"aws_accountid"`
	Alias 		htypes.String `tfsdk:"aws_alias"`
	Region 		htypes.String `tfsdk:"aws_accountregion"`

}

func (r *awsAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "AWS Account Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "CyberArk Privilege Cloud Credential ID- Generated from CyberArk after onboarding account into a safe.",
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Custom Account Name for customizing the object name in a safe.",
				Required: true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the Credential object.",
				Required: true,
			},
			"platform": schema.StringAttribute{
				Description: "Management Platform associated with the Database Credential.",
				Required: true,
			},
			"safe": schema.StringAttribute{
				Description: "Target Safe where the credential object will be onboarded.",
				Required: true,
			},
			"secrettype": schema.StringAttribute{
				Description: "Secret type of credentials, for AWS Accounts this value must be set to key.",
				Required: true,
			},
			"secret": schema.StringAttribute{
				Description: "Secret Key of the credential object.",
				Required: true,
				Sensitive: true,
			},
			"sm_manage": schema.BoolAttribute{
				Description: "Automatic Management of a credential. Optional Value.",
				Optional: true,
			},
			"sm_manage_reason": schema.StringAttribute{
				Description: "If sm_manage is false, provide reason why credential is not managed.",
				Optional: true,
			},
			"aws_kid": schema.StringAttribute{
				Description: "AWS Access Key ID.",
				Required: true,
			},
			"aws_accountid": schema.StringAttribute{
				Description: "AWS Account ID Number.",
				Required: true,
			},
			"aws_alias": schema.StringAttribute{
				Description: "AWS Account Alias.",
				Optional: true,
			},
			"aws_accountregion": schema.StringAttribute{
				Description: "AWS Region.",
				Optional: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *awsAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *awsAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan awsCredModel
	var props cybrtypes.AccountProps
	var sm_props cybrtypes.SecretManagement

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var name, username, platform, safe, secrettype, secret, aws_kid, aws_accountid, aws_alias, aws_accountregion, sm_manage_reason string

	var sm_manage bool

	name = os.Getenv("CYBERARK_ACCOUNT_CUSTOM_NAME")
	username = os.Getenv("CYBERARK_ACCOUNT_USERNAME")
	platform = os.Getenv("CYBERARK_ACCOUNT_PLATFORM")
	safe = os.Getenv("CYBERARK_ACCOUNT_SAFE")
	secrettype = os.Getenv("CYBERARK_ACCOUNT_SECRETTYPE")
	secret = os.Getenv("CYBERARK_ACCOUNT_SECRET")

	// Not processing env variable overrides for optional fields, this must be defined in the terraform plan.
	// This is a design decision.

	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
	}



	if !plan.Username.IsNull() {
		username = plan.Username.ValueString()
	}

	if !plan.Platform.IsNull() {
		platform = plan.Platform.ValueString()
	}

	if !plan.Platform.IsNull() {
		platform = plan.Platform.ValueString()
	}

	if !plan.Safe.IsNull() {
		safe = plan.Safe.ValueString()
	}

	if !plan.SecretType.IsNull() {
		secrettype = plan.SecretType.ValueString()
	}

	if !plan.Secret.IsNull() {
		secret = plan.Secret.ValueString()
	}
	
	if !plan.Manage.IsNull() {
		sm_manage = plan.Manage.ValueBool()
		sm_props.AutomaticManagement = &sm_manage
	}

	if !plan.ManageReason.IsNull() {
		sm_manage_reason = plan.ManageReason.ValueString()
		sm_props.ManualManagementReason = &sm_manage_reason
	}

	if !plan.AWSKID.IsNull() {
		aws_kid = plan.AWSKID.ValueString()
		props.AWSKID = &aws_kid
	}

	if !plan.AWSAccount.IsNull() {
		aws_accountid = plan.AWSAccount.ValueString()
		props.AWSAccount = &aws_accountid
	}

	if !plan.Alias.IsNull() {
		aws_alias = plan.Alias.ValueString()
		props.Alias = &aws_alias
	}

	if !plan.Region.IsNull() {
		aws_accountregion = plan.Region.ValueString()
		props.Region = &aws_accountregion
	}

	newAccount := cybrtypes.Credential {
		Name: &name,
		UserName: &username,
		Platform: &platform,
		SafeName: &safe,
		SecretType: &secrettype,
		Secret: &secret,
		Props: &props,
		SecretMgmt: &sm_props,
	}



	create, err := cybrapi.CreateAccount(&newAccount, r.client.AuthToken, r.client.Domain)
	if err != nil {
		log.Fatal(err)
	}

	if create == "" || len(create) == 0 {

		tflog.Error(ctx, "Error catching ID, this means provisioning failed in pipeline. Please check logs and try again.")

	} else {

		plan.ID = htypes.StringValue(create)
		plan.LastUpdated = htypes.StringValue(time.Now().Format(time.RFC850))
	
		// Set state to fully populated data
		resp.State.Set(ctx, plan)

	}
	
}

// Refresh Existing State
func (r *awsAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var currState awsCredModel
	diags := req.State.Get(ctx, &currState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := cybrapi.GetAccount(currState.ID.ValueStringPointer(), r.client.AuthToken, r.client.Domain)
	if err != nil {
		tflog.Error(ctx, "Unable to retrieve state from CyberArk API")
	}

	tflog.Info(ctx, "Refreshing state")

	// Main Credentials
	currState.Name = htypes.StringValue(*newState.Name)
	currState.Platform = htypes.StringValue(*newState.Platform)
	currState.Safe = htypes.StringValue(*newState.SafeName)
	currState.Username = htypes.StringValue(*newState.UserName)
	currState.SecretType = htypes.StringValue(*newState.SecretType)
	
	// AWS Props
	currState.AWSKID = htypes.StringValue(*newState.Props.AWSKID)
	currState.AWSAccount = htypes.StringValue(*newState.Props.AWSAccount)
	currState.Alias = htypes.StringValue(*newState.Props.Alias)
	currState.Region = htypes.StringValue(*newState.Props.Region)
	
	// Set last updated time to last updated tim in the vault
	newTime := time.Unix(*newState.SecretMgmt.ModifiedTime, 0)
	currState.LastUpdated = htypes.StringValue(newTime.Format(time.RFC850))

	// Ensure ID is consistent
	currState.ID = htypes.StringValue(*newState.CredID)

}


// Update updates the resource and sets the updated Terraform state on success.
func (r *awsAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Update is not supported through terraform. Please consult with your CyberArk Administrator to process account property updates.")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *awsAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	
		// Retrieve values from state
		var state awsCredModel
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	
		id := state.ID.ValueString()

		err := cybrapi.RemoveAccount(&id, r.client.AuthToken, r.client.Domain)
		if err != nil {
			tflog.Error(ctx, "Unable to remove account")
			return
		}

		tflog.Info(ctx, "Successfully removed account")
}