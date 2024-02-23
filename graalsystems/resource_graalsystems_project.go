package graalsystems

import (
	"context"

	sdk "github.com/graalsystems/sdk/go"
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
				Required:    true,
				Description: "The name of the project",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the project",
			},
			"banner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The banner of the project",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The labels of the project",
			},
			"locked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the project is locked or not",
			},
			"favorite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the project is a favorite project or not",
			},
			"badge": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The badge of the project",
			},
		},
	}
}

func resourceGraalSystemsProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	banner := d.Get("banner").(string)
	locked := d.Get("locked").(bool)
	favorite := d.Get("favorite").(bool)
	badge := d.Get("badge").(string)

	lbl := d.Get("labels").(map[string]interface{})
	labels := toStringMap(lbl)

	project := &sdk.Project{
		Name:        &name,
		Description: &description,
		Banner:      &banner,
		Labels:      &labels,
		Locked:      &locked,
		Favorite:    &favorite,
		Badge:       &badge,
	}
	result, _, err := apiClient.ProjectAPI.CreateProject(context.Background()).XTenant(meta.tenant).Project(*project).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	return resourceGraalSystemsProjectRead(ctx, d, meta)
}

func resourceGraalSystemsProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	project, _, err := apiClient.ProjectAPI.FindProjectById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", project.Name)
	if project.Description != nil {
		_ = d.Set("description", project.Description)
	}
	_ = d.Set("banner", project.Banner)
	_ = d.Set("labels", project.Labels)
	_ = d.Set("locked", project.Locked)
	_ = d.Set("favorite", project.Favorite)
	_ = d.Set("badge", project.Badge)

	return nil
}

func resourceGraalSystemsProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	if d.HasChange("name") {
		_, _, err := apiClient.ProjectAPI.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(patchFromResourceData(d, "name")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("description") {
		_, _, err := apiClient.ProjectAPI.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(patchFromResourceData(d, "description")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("banner") {
		_, _, err := apiClient.ProjectAPI.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(patchFromResourceData(d, "banner")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("badge") {
		_, _, err := apiClient.ProjectAPI.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(patchFromResourceData(d, "badge")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("labels") {
		_, _, err := apiClient.ProjectAPI.UpdateProject(context.Background(), d.Id()).XTenant(meta.tenant).Patch(patchFromResourceData(d, "labels")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// TODO add favorite & locked

	return resourceGraalSystemsProjectRead(ctx, d, meta)
}

func resourceGraalSystemsProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.ProjectAPI.DeleteProjectById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
