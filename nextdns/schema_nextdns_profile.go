package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Profile name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"endpoint": {
			Description: "Endpoints.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"doh": {
						Description: "The DNS over HTTPS address the profile is reachable at.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"dot": {
						Description: "The DNS over TLS address the profile is reachable at.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"ipv4": {
						Description: "The IPv4 addresses the profile is reachable at.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Computed: true,
					},
					"ipv6": {
						Description: "The IPv6 addresses the profile is reachable at.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Computed: true,
					},
				},
			},
		},
		"linkedip": {
			Description: "Linked IP.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"servers": {
						Description: "The DNS servers available for the profile.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Computed: true,
					},
					"ip": {
						Description: "The IP linked to the profile.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"update_token": {
						Description: "The update token to use to update the linked IP.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}
