package graalsystems

import (
	"context"
	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraalSystemsProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsProjectCreate,
		ReadContext:   resourceGraalSystemsProjectRead,
		UpdateContext: resourceGraalSystemsProjectUpdate,
		DeleteContext: resourceGraalSystemsProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the project",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the project",
			},
		},
	}
}

func resourceGraalSystemsProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	project := &sdk.Project{
		Name:        &name,
		Description: &description,
	}
	_, _, err := apiClient.ProjectApi.CreateProject(context.Background()).XTenant(meta.tenant).Project(*project).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGraalSystemsProjectRead(ctx, d, meta)
}

func resourceGraalSystemsProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.ProjectApi.FindProjectById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)

	return nil
}

func resourceGraalSystemsProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		_, _, err := apiClient.ProjectApi.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsProjectRead(ctx, d, meta)
}

func resourceGraalSystemsProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.ProjectApi.DeleteProjectById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
