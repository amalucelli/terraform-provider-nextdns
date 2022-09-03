package nextdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSSecuritySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Description: "The profile identifier to target the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"threat_intelligence_feeds": {
			Description: "Threat intelligence feeds.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"ai_threat_detection": {
			Description: "AI-Driven threat detection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"google_safe_browsing": {
			Description: "Google safe browsing.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"crypto_jacking": {
			Description: "Cryptojacking protection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"dns_rebinding": {
			Description: "DNS rebinding protection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"idn_homographs": {
			Description: "IDN homograph attacks protection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"typo_squatting": {
			Description: "Typosquatting protection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"dga": {
			Description: "Domain generation algorithms (DGAs) protection.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"nrd": {
			Description: "Block newly registered domains (NRDs).",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"ddns": {
			Description: "Block dynamic DNS hostnames.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"parking": {
			Description: "Block parked domains.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"csam": {
			Description: "Block child sexual abuse material.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"tlds": {
			Description: "Block top-level domains (TLDs).",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
