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
	}
}
