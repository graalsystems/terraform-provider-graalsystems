package graalsystems

import (
	"context"
	"fmt"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraalSystemsGroup() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceGraalSystemsGroup().Schema)

	dsSchema["group_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the group",
	}

	dsSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceGraalSystemsGroupRead,
		Schema:      dsSchema,
	}
}

func dataSourceGraalSystemsGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	var group sdk.Group
	if groupId := d.Get("group_id").(string); groupId != "" {
		p, _, err := apiClient.GroupAPI.FindGroupById(context.Background(), groupId).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		group = *p
	} else {
		groups, _, err := apiClient.GroupAPI.FindGroups(context.Background()).XTenant(meta.tenant).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		name := d.Get("name").(string)
		if name == "" {
			return diag.FromErr(fmt.Errorf("either group_id or name must be set"))
		}
		var matches []sdk.Group
		for _, item := range groups {
			if deref(item.Name) == name {
				matches = append(matches, item)
			}
		}
		groups = matches
		if len(groups) == 0 {
			return diag.FromErr(fmt.Errorf("no group found with the name %s", d.Get("name")))
		}
		if len(groups) > 1 {
			return diag.FromErr(fmt.Errorf("%d groups found with the same name %s", len(groups), d.Get("name")))
		}
		group = groups[0]
	}

	d.SetId(*group.Id)
	_ = d.Set("group_id", group.Id)

	return resourceGraalSystemsGroupRead(ctx, d, meta)
}
