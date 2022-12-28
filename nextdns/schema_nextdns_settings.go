package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"logs": {
			Description: "Logs.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Description: "Enable logs.",
						Type:        schema.TypeBool,
						Required:    true,
					},
					"privacy": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"log_clients_ip": {
									Description: "Log clients IP.",
									Type:        schema.TypeBool,
									Required:    true,
								},
								"log_domains": {
									Description: "Log domains.",
									Type:        schema.TypeBool,
									Required:    true,
								},
							},
						},
					},
					// TODO(amalucelli): Move this to string and parse to int in the provider.
					"retention": {
						Description: "Retention period for logs.",
						Type:        schema.TypeInt,
						Required:    true,
					},
					"location": {
						Description: "Location of the logs.",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
			},
		},
		"block_page": {
			Description: "Block Page.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Description: "Enable block page.",
						Type:        schema.TypeBool,
						Required:    true,
					},
				},
			},
		},
		"performance": {
			Description: "Performance.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ecs": {
						Description: "Anonymized EDNS Client Subnet.",
						Type:        schema.TypeBool,
						Required:    true,
					},
					"cache_boost": {
						Description: "Cache Boost.",
						Type:        schema.TypeBool,
						Required:    true,
					},
					"cname_flattening": {
						Description: "CNAME Flattening.",
						Type:        schema.TypeBool,
						Required:    true,
					},
				},
			},
		},
		"web3": {
			Description: "Web3.",
			Type:        schema.TypeBool,
			Required:    true,
		},
	}
}
