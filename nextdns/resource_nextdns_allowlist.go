package nextdns

import (
	"context"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceNextDNSAllowlist() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSAllowlistSchema(),
		CreateContext: resourceNextDNSAllowlistCreate,
		ReadContext:   resourceNextDNSAllowlistRead,
		UpdateContext: resourceNextDNSAllowlistUpdate,
		DeleteContext: resourceNextDNSAllowlistDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSAllowlistImport,
		},
	}
}

func resourceNextDNSAllowlistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	if domain, ok := d.GetOk("domain"); ok {
		var Allowlist []*nextdns.Allowlist
		for _, elem := range domain.(*schema.Set).List() {
			Allowlist = append(Allowlist, &nextdns.Allowlist{
				ID:     elem.(map[string]interface{})["id"].(string),
				Active: elem.(map[string]interface{})["active"].(bool),
			})
		}

		request := &nextdns.CreateAllowlistRequest{
			ProfileID: profileID,
			Allowlist: Allowlist,
		}
		err := client.Allowlist.Create(ctx, request)
		if err != nil {
			return diag.FromErr(errors.Wrap(err, "error creating allow list"))
		}
	}

	d.SetId(profileID)

	return resourceNextDNSAllowlistRead(ctx, d, meta)
}

func resourceNextDNSAllowlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetAllowlistRequest{
		ProfileID: profileID,
	}
	Allowlist, err := client.Allowlist.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting allow list"))
	}

	var domains []map[string]interface{}
	var domain map[string]interface{}

	for _, d := range Allowlist {
		domain = make(map[string]interface{})
		domain["id"] = d.ID
		domain["active"] = d.Active

		domains = append(domains, domain)
	}
	if err := d.Set("domain", domains); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(profileID)

	return nil
}

func resourceNextDNSAllowlistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	if domain, ok := d.GetOk("domain"); ok {
		var Allowlist []*nextdns.Allowlist
		for _, elem := range domain.(*schema.Set).List() {
			Allowlist = append(Allowlist, &nextdns.Allowlist{
				ID:     elem.(map[string]interface{})["id"].(string),
				Active: elem.(map[string]interface{})["active"].(bool),
			})
		}

		request := &nextdns.CreateAllowlistRequest{
			ProfileID: profileID,
			Allowlist: Allowlist,
		}
		err := client.Allowlist.Create(ctx, request)
		if err != nil {
			return diag.FromErr(errors.Wrap(err, "error creating allow list"))
		}
	}

	return resourceNextDNSAllowlistRead(ctx, d, meta)
}

func resourceNextDNSAllowlistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.CreateAllowlistRequest{
		ProfileID: profileID,
		Allowlist: []*nextdns.Allowlist{},
	}
	err := client.Allowlist.Create(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting allow list"))
	}

	return resourceNextDNSAllowlistRead(ctx, d, meta)
}

func resourceNextDNSAllowlistImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSAllowlistRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
