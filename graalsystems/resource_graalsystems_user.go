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
				Description: "The username of the user",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email address of the user",
			},
			"first_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first name of the user",
			},
			"last_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The last name of the user",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: "Initial password of the user, set at creation only. Without it, the user " +
					"signs in through an identity provider or after a password reset.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "The platform does not store a description for users; this attribute is ignored.",
				Description: "Ignored: the platform does not store a description for users.",
			},
		},
	}
}

func resourceGraalSystemsUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	username := d.Get("name").(string)
	email := d.Get("email").(string)
	user := &sdk.User{
		Username: &username,
		Email:    &email,
	}
	if v, ok := d.GetOk("first_name"); ok {
		s := v.(string)
		user.Firstname = &s
	}
	if v, ok := d.GetOk("last_name"); ok {
		s := v.(string)
		user.Lastname = &s
	}
	if v, ok := d.GetOk("password"); ok {
		s := v.(string)
		user.Password = &s
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

	// The password is never returned, and the description is not stored: both keep
	// the configured value.
	_ = d.Set("name", res.Username)
	_ = d.Set("email", res.Email)
	_ = d.Set("first_name", res.Firstname)
	_ = d.Set("last_name", res.Lastname)

	return nil
}

func resourceGraalSystemsUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)

	if d.HasChange("password") {
		return diag.Errorf("password is set at creation only: change it in the identity server " +
			"(password reset), or recreate the user")
	}

	// Terraform attribute → API field.
	var ops []patchOperation
	for attr, field := range map[string]string{
		"name":       "username",
		"email":      "email",
		"first_name": "firstname",
		"last_name":  "lastname",
	} {
		if d.HasChange(attr) {
			ops = append(ops, patchOperation{Op: "add", Path: "/" + field, Value: d.Get(attr)})
		}
	}
	if err := applyPatch(ctx, meta, "user", "/users/"+d.Id(), ops); err != nil {
		return diag.FromErr(err)
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
