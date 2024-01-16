package graalsystems

import (
	"context"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraalSystemsUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsUserCreate,
		ReadContext:   resourceGraalSystemsUserRead,
		UpdateContext: resourceGraalSystemsUserUpdate,
		DeleteContext: resourceGraalSystemsUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the user",
			},
		},
	}
}

func resourceGraalSystemsUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	username := d.Get("username").(string)
	description := d.Get("description").(string)
	user := &sdk.User{
		Username:    &username,
		Description: &description,
	}
	result, _, err := apiClient.UserAPI.CreateUser(context.Background()).XTenant(meta.tenant).User(*user).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)

	return resourceGraalSystemsUserRead(ctx, d, meta)
}

func resourceGraalSystemsUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.UserAPI.FindUserById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("username", res.Username)
	_ = d.Set("description", res.Description)

	return nil
}

func resourceGraalSystemsUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	if d.HasChange("username") {
		path := "/username"

		value := d.Get("username").(string)

		patch := &sdk.Patch{
			Op:    nil,
			Path:  &path,
			Value: &value,
		}
		patchs := &[]sdk.Patch{*patch}
		_, _, err := apiClient.UserAPI.UpdateUser(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsUserRead(ctx, d, meta)
}

func resourceGraalSystemsUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.UserAPI.DeleteUserById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
