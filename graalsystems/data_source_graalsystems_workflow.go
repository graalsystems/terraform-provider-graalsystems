package graalsystems

import (
	"context"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

// dataSourceGraalSystemsWorkflow returns a datasource that can be used to retrieve a workflow from the GraalSystems API
func dataSourceGraalSystemsWorkflow() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsWorkflow().Schema)

	dsSchema["workflow_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the workflow",
	}
	dsSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsWorkflowRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsWorkflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	// Retrieve the input
	workflowId := d.Get("workflow_id").(string)
	name := d.Get("name").(string)

	if workflowId == "" && name == "" {
		return diag.FromErr(fmt.Errorf("workflow_id or name must be set"))
	}
	if workflowId != "" && name != "" {
		return diag.FromErr(fmt.Errorf("workflow_id and name cannot be set at the same time"))
	}

	var filteredWorkflow *sdk.Workflow
	// Retrieving the workflow by its id is straightforward
	if workflowId != "" {
		if res, _, err := apiClient.WorkflowAPI.FindWorkflowById(context.Background(), workflowId).XTenant(meta.tenant).Execute(); err != nil {
			diag.FromErr(err)
		} else {
			filteredWorkflow = res
		}
	}
	// Retrieving the workflow by its name need to retrieve all the workflows and filter them
	if name != "" {
		if page, _, err := apiClient.WorkflowAPI.FindWorkflows(context.Background()).XTenant(meta.tenant).Execute(); err != nil {
			return diag.FromErr(err)
		} else {
			var matches []sdk.Workflow

			for _, workflow := range page.Content {
				if strings.TrimSpace(*workflow.Name) == strings.TrimSpace(name) {
					matches = append(matches, workflow)
				}
			}
			if len(matches) == 0 {
				return diag.FromErr(fmt.Errorf("no workflow exists with the name %s", name))
			}
			if len(matches) > 1 {
				return diag.FromErr(fmt.Errorf("%d workflows exist with the same name %s. You can filter them by their id", len(matches), name))
			}
			filteredWorkflow = &matches[0]
			// Retrieving additional information about the workspace
			if workflow, _, err := apiClient.WorkflowAPI.FindWorkflowById(context.Background(), *filteredWorkflow.Id).XTenant(meta.tenant).Execute(); err != nil {
				return diag.FromErr(err)
			} else {
				filteredWorkflow = workflow
			}
		}
	}
	d.SetId(*filteredWorkflow.Id)
	return resourceGraalSystemsWorkflowRead(ctx, d, m)
}
