package graalsystems

import (
	"context"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraalSystemsGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsGroupCreate,
		ReadContext:   resourceGraalSystemsGroupRead,
		UpdateContext: resourceGraalSystemsGroupUpdate,
		DeleteContext: resourceGraalSystemsGroupDelete,
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
		},
	}
}

func resourceGraalSystemsGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	project := &sdk.Group{
		Name:        &name,
		Description: &description,
	}
	result, _, err := apiClient.GroupAPI.CreateGroup(context.Background()).XTenant(meta.tenant).Group(*project).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	return resourceGraalSystemsGroupRead(ctx, d, meta)
}

func resourceGraalSystemsGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.GroupAPI.FindGroupById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
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

func resourceGraalSystemsGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	//apiClient := meta.apiClient

	//if d.HasChange("name") {
	//	path := "/name"
	//
	//	value := make(map[string]interface{})
	//	value["name"] = d.Get("name").(string)
	//
	//	patch := &sdk.Patch{
	//		Op:    nil,
	//		Path:  &path,
	//		Value: &value,
	//	}
	//	patchs := &[]sdk.Patch{*patch}
	//	_, _, err := apiClient.GroupAPI.UpdateGroup(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//}

	return resourceGraalSystemsGroupRead(ctx, d, meta)
}

func resourceGraalSystemsGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.GroupAPI.DeleteGroupById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
