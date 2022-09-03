package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSPrivacySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"allow_affiliate": {
			Description: "Allow affiliate & tracking links.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"disguised_trackers": {
			Description: "Block disguised third-party trackers.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"blocklists": {
			Description: "Blocklists.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"natives": {
			Description: "Native tracking protection.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
