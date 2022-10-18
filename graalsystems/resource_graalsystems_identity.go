package graalsystems

import (
	"context"
	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraalSystemsIdentity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsIdentityCreate,
		ReadContext:   resourceGraalSystemsIdentityRead,
		UpdateContext: resourceGraalSystemsIdentityUpdate,
		DeleteContext: resourceGraalSystemsIdentityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the identity",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the identity",
			},
		},
	}
}

func resourceGraalSystemsIdentityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	identity := &sdk.Identity{
		Name:        &name,
		Description: &description,
	}
	_, _, err := apiClient.IdentityApi.CreateIdentity(context.Background()).XTenant(meta.tenant).Identity(*identity).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGraalSystemsIdentityRead(ctx, d, meta)
}

func resourceGraalSystemsIdentityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.IdentityApi.FindIdentityById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
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

func resourceGraalSystemsIdentityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		_, _, err := apiClient.IdentityApi.UpdateIdentity(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsIdentityRead(ctx, d, meta)
}

func resourceGraalSystemsIdentityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.IdentityApi.DeleteIdentityById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
