package nextdns

import (
	"context"
	"os"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NextDNS API Key",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nextdns_setup_endpoint": dataSourceNextDNSSetupEndpoint(),
			"nextdns_setup_linkedip": dataSourceNextDNSSetupLinkedIP(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nextdns_allowlist":        resourceNextDNSAllowlist(),
			"nextdns_denylist":         resourceNextDNSDenylist(),
			"nextdns_parental_control": resourceNextDNSParentalControl(),
			"nextdns_privacy":          resourceNextDNSPrivacy(),
			"nextdns_profile":          resourceNextDNSProfile(),
			"nextdns_rewrite":          resourceNextDNSRewrite(),
			"nextdns_security":         resourceNextDNSSecurity(),
			"nextdns_settings":         resourceNextDNSSettings(),
		},
		ConfigureContextFunc: configure,
	}
}

// nolint:revive
func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := os.Getenv("NEXTDNS_API_KEY")

	if key, ok := d.Get("api_key").(string); ok && len(key) > 0 {
		apiKey = key
	}

	if len(apiKey) == 0 {
		return nil, diag.Errorf(
			"NextDNS API key must be provided in the provider block or NEXTDNS_API_KEY environment variable.",
		)
	}

	client, err := nextdns.New(nextdns.WithAPIKey(apiKey))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
