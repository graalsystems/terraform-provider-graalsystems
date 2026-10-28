package graalsystems

import (
	"context"
	"fmt"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsProject() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsProject().Schema)

	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the project",
	}

	dsSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsProjectRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var project sdk.Project
	if projectId := d.Get("project_id").(string); projectId != "" {
		p, _, err := apiClient.ProjectAPI.FindProjectById(context.Background(), projectId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		project = *p
	} else {
		projects, _, err := apiClient.ProjectAPI.FindProjects(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		name := d.Get("name").(string)
		if name == "" {
			return diag.FromErr(fmt.Errorf("either project_id or name must be set"))
		}
		var matches []sdk.Project
		for _, item := range projects {
			if deref(item.Name) == name {
				matches = append(matches, item)
			}
		}
		projects = matches
		if len(projects) == 0 {
			return diag.FromErr(fmt.Errorf("no project found with the name %s", d.Get("name")))
		}
		if len(projects) > 1 {
			return diag.FromErr(fmt.Errorf("%d projects found with the same name %s", len(projects), d.Get("name")))
		}
		project = projects[0]
	}

	d.SetId(*project.Id)
	_ = d.Set("project_id", project.Id)

	return resourceGraalSystemsProjectRead(ctx, d, meta)
}
