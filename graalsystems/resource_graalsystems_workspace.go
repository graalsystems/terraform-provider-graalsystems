package graalsystems

import (
	"context"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"slices"
)

const (
	workspaceTypeJupyter  = "jupyter"
	workspaceTypeMetabase = "metabase"
	workspaceTypeSuperset = "superset"
	workspaceTypeVsCode   = "vscode"
	workspaceTypeZeppelin = "zeppelin"
)

var workspaceTypes = []string{workspaceTypeJupyter, workspaceTypeMetabase, workspaceTypeSuperset, workspaceTypeVsCode, workspaceTypeZeppelin}

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
				Description: "The type of the workspace",
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
				Description: "The id of the infrastructure to deploy the workspace on",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance type of the compute used for the workspace",
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
	if result, request, err := apiClient.WorkspaceAPI.CreateWorkspace(context.Background()).XTenant(meta.tenant).Workspace(*workspace).Execute(); err != nil {
		return diag.FromErr(err)
	} else if request != nil && request.StatusCode == 200 {
		return diag.FromErr(fmt.Errorf("workspace created, but could not retrieve its info. Check that every parameter you entered is valid. InfrastructureId:%s ; InstanceType:%s", infrastructureId, instanceType))
	} else {
		d.SetId(*result.Id)
	}

	return resourceGraalSystemsWorkspaceRead(ctx, d, meta)
}

func resourceGraalSystemsWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	_ = d.Set("owner", *res.Owner)
	_ = d.Set("version", *res.Version)
	_ = d.Set("status", *res.Status)
	_ = d.Set("public_url", *res.PublicUrl)

	return nil
}

func patchFromResourceData(d *schema.ResourceData, patchElement string) *sdk.Patch {
	path := "/" + patchElement
	value := make(map[string]interface{})
	value[patchElement] = d.Get(patchElement).(string)
	operation := "replace"
	return &sdk.Patch{Op: &operation, Path: &path, Value: value}
}

func resourceGraalSystemsWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient
	//TODO: fix this deserialization issue

	if d.HasChange("name") {
		_, _, err := apiClient.WorkspaceAPI.UpdateWorkspace(context.Background(), d.Id()).XTenant(meta.tenant).Patch([]sdk.Patch{*patchFromResourceData(d, "name")}).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		_, _, err := apiClient.WorkspaceAPI.UpdateWorkspace(context.Background(), d.Id()).XTenant(meta.tenant).Patch([]sdk.Patch{*patchFromResourceData(d, "description")}).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("type") {
		_, _, err := apiClient.WorkspaceAPI.UpdateWorkspace(context.Background(), d.Id()).XTenant(meta.tenant).Patch([]sdk.Patch{*patchFromResourceData(d, "type")}).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("infrastructure_id") {
		_, _, err := apiClient.WorkspaceAPI.UpdateWorkspace(context.Background(), d.Id()).XTenant(meta.tenant).Patch([]sdk.Patch{*patchFromResourceData(d, "infrastructure_id")}).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("instance_type") {
		_, _, err := apiClient.WorkspaceAPI.UpdateWorkspace(context.Background(), d.Id()).XTenant(meta.tenant).Patch([]sdk.Patch{*patchFromResourceData(d, "instance_type")}).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsJobRead(ctx, d, meta)
}

func resourceGraalSystemsWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.WorkspaceAPI.DeleteWorkspaceById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
