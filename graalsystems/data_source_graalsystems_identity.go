package graalsystems

import (
	"context"
	"fmt"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsIdentity() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsIdentity().Schema)

	dsSchema["identity_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the identity",
	}

	dsSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsIdentityRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsIdentityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var identity sdk.Identity
	if identityId := d.Get("identity_id").(string); identityId != "" {
		p, _, err := apiClient.IdentityAPI.FindIdentityById(context.Background(), identityId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		identity = *p
	} else {
		identities, _, err := apiClient.IdentityAPI.FindIdentities(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		name := d.Get("name").(string)
		if name == "" {
			return diag.FromErr(fmt.Errorf("either identity_id or name must be set"))
		}
		var matches []sdk.Identity
		for _, item := range identities {
			if deref(item.Name) == name {
				matches = append(matches, item)
			}
		}
		identities = matches
		if len(identities) == 0 {
			return diag.FromErr(fmt.Errorf("no identity found with the name %s", d.Get("name")))
		}
		if len(identities) > 1 {
			return diag.FromErr(fmt.Errorf("%d identities found with the same name %s", len(identities), d.Get("name")))
		}
		identity = identities[0]
	}

	d.SetId(*identity.Id)
	_ = d.Set("identity_id", identity.Id)

	return resourceGraalSystemsIdentityRead(ctx, d, meta)
}
