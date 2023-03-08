package nextdns

import (
	"context"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceNextDNSSetupLinkedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNextDNSSetupLinkedIPRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Description: "The profile identifier to target the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
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
			"ddns": {
				Description: "The DDNS configuration for the linked IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"update_token": {
				Description: "The update token to use to update the linked IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceNextDNSSetupLinkedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetSetupLinkedIPRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	setup, err := client.SetupLinkedIP.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting setup linkedip settings"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", setup))

	d.SetId(profileID)
	d.Set("servers", setup.Servers)
	d.Set("ip", setup.IP)
	d.Set("ddns", setup.Ddns)
	d.Set("update_token", setup.UpdateToken)

	return nil
}
