package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSParentalControlSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"block_bypass": {
			Description: "Block bypass methods.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"category": {
			Description: "Restrict access to specific categories of websites and apps.",
			Type:        schema.TypeSet,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"active": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"safe_search": {
			Description: "Safe search.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"service": {
			Description: "Restrict access to specific websites, apps and games.",
			Type:        schema.TypeSet,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"active": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"youtube_restricted_mode": {
			Description: "YouTube restricted mode.",
			Type:        schema.TypeBool,
			Required:    true,
		},
	}
}
