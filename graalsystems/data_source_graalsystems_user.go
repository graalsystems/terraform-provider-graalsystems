package graalsystems

import (
	"context"
	"fmt"

	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsUser() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsUser().Schema)

	dsSchema["user_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the user",
	}

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsUserRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var user sdk.User
	userId, ok := d.Get("user_id").(string)
	if ok {
		p, _, err := apiClient.UserApi.FindUserById(context.Background(), userId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		user = p
	} else {
		users, _, err := apiClient.UserApi.FindUsers(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(users) == 0 {
			return diag.FromErr(fmt.Errorf("no user found with the name %s", d.Get("name")))
		}
		if len(users) > 1 {
			return diag.FromErr(fmt.Errorf("%d users found with the same name %s", len(users), d.Get("name")))
		}
		user = users[0]
	}

	d.SetId(*user.Id)
	_ = d.Set("user_id", user.Id)

	return resourceGraalSystemsUserRead(ctx, d, meta)
}
