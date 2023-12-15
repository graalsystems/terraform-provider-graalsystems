package graalsystems

import (
	"context"
	"fmt"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsJob() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsJob().Schema)

	dsSchema["job_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the job",
	}

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsJobRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var job sdk.Job
	jobId, ok := d.Get("job_id").(string)
	if ok {
		p, _, err := apiClient.JobAPI.FindJobByJobId(context.Background(), jobId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		job = *p
	} else {
		jobs, _, err := apiClient.JobAPI.FindJobs(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(jobs) == 0 {
			return diag.FromErr(fmt.Errorf("no job found with the name %s", d.Get("name")))
		}
		if len(jobs) > 1 {
			return diag.FromErr(fmt.Errorf("%d jobs found with the same name %s", len(jobs), d.Get("name")))
		}
		job = jobs[0]
	}

	d.SetId(*job.Id)
	_ = d.Set("job_id", job.Id)

	return resourceGraalSystemsJobRead(ctx, d, meta)
}
