package graalsystems

import (
	"context"
	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraalSystemsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsJobCreate,
		ReadContext:   resourceGraalSystemsJobRead,
		UpdateContext: resourceGraalSystemsJobUpdate,
		DeleteContext: resourceGraalSystemsJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the job",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the job",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the job",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout of the job",
			},
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the job",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The labels associated with the job",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spark": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "List of private network to connect with your instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main_class_name": {
							Type:        schema.TypeString,
							Required:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"tensorflow": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "List of private network to connect with your instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number_replicas": {
							Type:        schema.TypeInt,
							Required:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceGraalSystemsJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	projectId := d.Get("project_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	job := &sdk.Job{
		Name:        &name,
		Description: &description,
	}
	result, _, err := apiClient.ProjectApi.CreateJobForProject(context.Background(), projectId).XTenant(meta.tenant).Job(*job).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	return resourceGraalSystemsJobRead(ctx, d, meta)
}

func resourceGraalSystemsJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.JobApi.FindJobByJobId(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("description", res.Description)

	return nil
}

func resourceGraalSystemsJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	if d.HasChange("name") {
		path := "/name"

		value := make(map[string]interface{})
		value["name"] = d.Get("name").(string)

		patch := &sdk.Patch{
			Op:    nil,
			Path:  &path,
			Value: &value,
		}
		patchs := &[]sdk.Patch{*patch}
		_, _, err := apiClient.JobApi.UpdateJob(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsJobRead(ctx, d, meta)
}

func resourceGraalSystemsJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.JobApi.DeleteJobById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
