package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSRewriteSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"rewrite": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"domain": {
						Type:     schema.TypeString,
						Required: true,
					},
					"address": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}
