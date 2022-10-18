package graalsystems

import (
	"context"
	"fmt"

	"github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsIdentity() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsIdentity().Schema)

	dsSchema["name"].ConflictsWith = []string{"identity_id"}
	dsSchema["identity_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the identity",
	}

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsIdentityRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsIdentityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var identity sdk.Identity
	identityId, ok := d.Get("identity_id").(string)
	if ok {
		p, _, err := apiClient.IdentityApi.FindIdentityById(context.Background(), identityId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		identity = p
	} else {
		identitys, _, err := apiClient.IdentityApi.FindIdentities(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(identitys) == 0 {
			return diag.FromErr(fmt.Errorf("no identity found with the name %s", d.Get("name")))
		}
		if len(identitys) > 1 {
			return diag.FromErr(fmt.Errorf("%d identitys found with the same name %s", len(identitys), d.Get("name")))
		}
		identity = identitys[0]
	}

	d.SetId(*identity.Id)
	_ = d.Set("identity_id", identity.Id)

	return resourceGraalSystemsIdentityRead(ctx, d, meta)
}
