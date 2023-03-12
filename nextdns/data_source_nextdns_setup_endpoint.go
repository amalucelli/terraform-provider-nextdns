package nextdns

import (
	"context"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNextDNSSetupEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNextDNSSetupEndpointRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Description: "The profile identifier to target the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
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
			"dnscrypt": {
				Description: "The DNS Stamps from the profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceNextDNSSetupEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetSetupRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	setup, err := client.Setup.Get(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting setup endpoint settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", setup))

	d.SetId(profileID)
	d.Set("doh", DNSOverHTTPSAddress(profileID))
	d.Set("dot", DNSOverTLSAddress(profileID))
	d.Set("ipv4", setup.Ipv4)
	d.Set("ipv6", setup.Ipv6)
	d.Set("dnscrypt", setup.Dnscrypt)

	return nil
}
