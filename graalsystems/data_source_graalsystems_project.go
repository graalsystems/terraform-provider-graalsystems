package graalsystems

import (
	"context"
	"fmt"

	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsProject() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsProject().Schema)

	dsSchema["name"].ConflictsWith = []string{"project_id"}
	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the project",
	}

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsProjectRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var project sdk.Project
	projectId, ok := d.Get("project_id").(string)
	if ok {
		p, _, err := apiClient.ProjectApi.FindProjectById(context.Background(), projectId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		project = p
	} else {
		projects, _, err := apiClient.ProjectApi.FindProjects(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
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
