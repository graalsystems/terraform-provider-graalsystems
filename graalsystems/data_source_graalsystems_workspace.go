package graalsystems

import (
	"context"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceGraalSystemsWorkspace returns a datasource that can be used to retrieve a workspace from the GraalSystems API
func dataSourceGraalSystemsWorkspace() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsWorkspace().Schema)

	dsSchema["workspace_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the workspace",
	}
	dsSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsWorkspaceRead,
		Schema:      dsSchema,
	}
}

// dataSourceGraalSystemsWorkspaceRead reads the workspace from the GraalSystems API and returns its attributes
// The workspace can be retrieved by its id or its name
func dataSourceGraalSystemsWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	// Retrieve the input
	workspaceId := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	if workspaceId == "" && name == "" {
		return diag.FromErr(fmt.Errorf("workspace_id or name must be set"))
	}
	if workspaceId != "" && name != "" {
		return diag.FromErr(fmt.Errorf("workspace_id and name cannot be set at the same time"))
	}

	var filteredWorkspace *sdk.Workspace
	// Retrieving the workspace by its id is straightforward
	if workspaceId != "" {
		if res, _, err := apiClient.WorkspaceAPI.FindWorkspaceById(context.Background(), workspaceId).XTenant(meta.tenant).Execute(); err != nil {
			diag.FromErr(err)
		} else {
			filteredWorkspace = res
		}
	}
	// Retrieving the workspace by its name need to retrieve all the workspaces and filter them
	if name != "" {
		if page, _, err := apiClient.WorkspaceAPI.FindWorkspaces(context.Background()).XTenant(meta.tenant).Execute(); err != nil {
			return diag.FromErr(err)
		} else {
			var matches []sdk.Workspace

			for _, space := range page.Content {
				if strings.TrimSpace(*space.Name) == strings.TrimSpace(name) {
					matches = append(matches, space)
				}
			}
			if len(matches) == 0 {
				return diag.FromErr(fmt.Errorf("no workspace exists with the name %s", name))
			}
			if len(matches) > 1 {
				return diag.FromErr(fmt.Errorf("%d workspaces exist with the same name %s. You can filter them by their id", len(matches), name))
			}
			filteredWorkspace = &matches[0]
			// Retrieving additional information about the workspace
			if space, _, err := apiClient.WorkspaceAPI.FindWorkspaceById(context.Background(), *filteredWorkspace.Id).XTenant(meta.tenant).Execute(); err != nil {
				return diag.FromErr(err)
			} else {
				filteredWorkspace = space
			}
		}
	}

	d.SetId(*filteredWorkspace.Id)
	_ = d.Set("name", filteredWorkspace.Name)
	if filteredWorkspace.Description != nil {
		_ = d.Set("description", *filteredWorkspace.Description)
	}
	_ = d.Set("type", *filteredWorkspace.Type)
	_ = d.Set("infrastructure_id", *filteredWorkspace.InfrastructureId)
	_ = d.Set("instance_type", *filteredWorkspace.InstanceType)
	_ = d.Set("owner", *filteredWorkspace.Owner)
	_ = d.Set("version", *filteredWorkspace.Version)
	_ = d.Set("status", *filteredWorkspace.Status)
	_ = d.Set("public_url", *filteredWorkspace.PublicUrl)

	return nil
}
