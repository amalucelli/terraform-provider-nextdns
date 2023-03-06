package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint_doh": {
			Description: "The DNS over HTTPS address the DNS is reachable at.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"endpoint_dot": {
			Description: "The DNS over TLS address the DNS is reachable at.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"endpoint_ipv6": {
			Description: "The IPv6 addresses the DNS is reachable at.",
			Type:        schema.TypeList,
			Elem:        schema.TypeString,
			Computed:    true,
		},
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
	}
}
