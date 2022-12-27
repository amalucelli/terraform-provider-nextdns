package nextdns

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
					"recreation": {
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
					"recreation": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"recreation": {
			Description: "Period for each day of the week during which some of the websites, apps, games or categories will not be blocked.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"timezone": {
						Type:     schema.TypeString,
						Required: true,
					},
					"monday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"tuesday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"wednesday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"thursday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"friday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"saturday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
					},
					"sunday": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     recreationTimeElem,
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

var recreationTimeElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"start": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]:00$`), "Must be in HH:MM:00 format"),
		},
		"end": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]:00$`), "Must be in HH:MM:00 format"),
		},
	},
}
