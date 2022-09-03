package nextdns

import (
	"context"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "NextDNS API Key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nextdns_denylist":  resourceNextDNSDenylist(),
			"nextdns_allowlist": resourceNextDNSAllowlist(),
			"nextdns_security":  resourceNextDNSSecurity(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client, err := nextdns.New(
		nextdns.WithAPIKey(d.Get("api_key").(string)))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
