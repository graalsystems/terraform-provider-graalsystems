package graalsystems

import (
	"context"
	"fmt"
	"slices"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Workspace types served by the platform.
const (
	workspaceTypeJupyter  = "jupyter"
	workspaceTypeVsCode   = "vscode"
	workspaceTypeRStudio  = "rstudio"
	workspaceTypeGrafana  = "grafana"
	workspaceTypeSuperset = "superset"
	workspaceTypeMetabase = "metabase"
	workspaceTypeZeppelin = "zeppelin"
)

var workspaceTypes = []string{
	workspaceTypeJupyter, workspaceTypeVsCode, workspaceTypeRStudio,
	workspaceTypeGrafana, workspaceTypeSuperset, workspaceTypeMetabase, workspaceTypeZeppelin,
}

// resourceGraalSystemsWorkspace defines the schema for the workspace resource
// It also implements the CRUD methods
func resourceGraalSystemsWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsWorkspaceCreate,
		ReadContext:   resourceGraalSystemsWorkspaceRead,
		UpdateContext: resourceGraalSystemsWorkspaceUpdate,
		DeleteContext: resourceGraalSystemsWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the workspace",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the workspace",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the workspace. Changing it recreates the workspace",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !slices.Contains(workspaceTypes, val.(string)) {
						errs = append(errs, fmt.Errorf("%q must be one of %q, got: %s", key, workspaceTypes, val))
					}
					return
				},
			},
			"infrastructure_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the infrastructure to deploy the workspace on. Changing it recreates the workspace",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance type of the compute used for the workspace. Changing it recreates the workspace",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the workspace",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the workspace",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the workspace type",
			},
			"public_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public url of the workspace",
			},
		},
	}
}

// resourceGraalSystemsWorkspaceCreate creates a workspace
func resourceGraalSystemsWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	user, _, err := apiClient.UserAPI.FindCurrentUser(context.Background()).XTenant(meta.tenant).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	workspaceType := d.Get("type").(string)
	infrastructureId := d.Get("infrastructure_id").(string)
	instanceType := d.Get("instance_type").(string)
	workspace := &sdk.Workspace{
		Name:             &name,
		Description:      &description,
		Type:             &workspaceType,
		InfrastructureId: &infrastructureId,
		InstanceType:     &instanceType,
		Owner:            user.Id,
	}
	// The API answers 201 (or 200, both declared by the contract) with the workspace. A 200 used to be
	// reported as an error, which left the created workspace outside the state.
	result, _, err := apiClient.WorkspaceAPI.CreateWorkspace(context.Background()).XTenant(meta.tenant).Workspace(*workspace).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil || result.Id == nil {
		return diag.FromErr(fmt.Errorf("workspace created, but the API returned no id (infrastructure_id %s, instance_type %s)", infrastructureId, instanceType))
	}
	d.SetId(*result.Id)

	return resourceGraalSystemsWorkspaceRead(ctx, d, meta)
}

// resourceGraalSystemsWorkspaceRead reads a workspace
func resourceGraalSystemsWorkspaceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.WorkspaceAPI.FindWorkspaceById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	if res.Description != nil {
		_ = d.Set("description", res.Description)
	}
	_ = d.Set("type", res.Type)
	_ = d.Set("infrastructure_id", res.InfrastructureId)
	_ = d.Set("instance_type", res.InstanceType)
	// Version, status and public URL are empty while the workspace starts: set the pointers as they
	// are (nil becomes ""), never dereference them.
	_ = d.Set("owner", res.Owner)
	_ = d.Set("version", res.Version)
	_ = d.Set("status", res.Status)
	_ = d.Set("public_url", res.PublicUrl)

	return nil
}

// resourceGraalSystemsWorkspaceUpdate updates a workspace
func resourceGraalSystemsWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)

	ops := replaceChanged(d, "name", "description")
	if err := applyPatch(ctx, meta, "workspace", "/workspaces/"+d.Id(), ops); err != nil {
		return diag.FromErr(err)
	}

	return resourceGraalSystemsWorkspaceRead(ctx, d, meta)
}

// resourceGraalSystemsWorkspaceDelete deletes a workspace
func resourceGraalSystemsWorkspaceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.WorkspaceAPI.DeleteWorkspaceById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
