package nextdns

import (
	"context"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceNextDNSDenylist() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSDenylistSchema(),
		CreateContext: resourceNextDNSDenylistCreate,
		ReadContext:   resourceNextDNSDenylistRead,
		UpdateContext: resourceNextDNSDenylistUpdate,
		DeleteContext: resourceNextDNSDenylistDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSDenylistImport,
		},
	}
}

func resourceNextDNSDenylistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	if domain, ok := d.GetOk("domain"); ok {
		var denylist []*nextdns.Denylist
		for _, elem := range domain.(*schema.Set).List() {
			denylist = append(denylist, &nextdns.Denylist{
				ID:     elem.(map[string]interface{})["id"].(string),
				Active: elem.(map[string]interface{})["active"].(bool),
			})
		}

		request := &nextdns.CreateDenylistRequest{
			ProfileID: profileID,
			Denylist:  denylist,
		}
		err := client.Denylist.Create(ctx, request)
		if err != nil {
			return diag.FromErr(errors.Wrap(err, "error creating deny list"))
		}
	}

	d.SetId(profileID)

	return resourceNextDNSDenylistRead(ctx, d, meta)
}

func resourceNextDNSDenylistRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetDenylistRequest{
		ProfileID: profileID,
	}
	denylist, err := client.Denylist.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting deny list"))
	}

	var domains []map[string]interface{}
	var domain map[string]interface{}

	for _, d := range denylist {
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

func resourceNextDNSDenylistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	if domain, ok := d.GetOk("domain"); ok {
		var denylist []*nextdns.Denylist
		for _, elem := range domain.(*schema.Set).List() {
			denylist = append(denylist, &nextdns.Denylist{
				ID:     elem.(map[string]interface{})["id"].(string),
				Active: elem.(map[string]interface{})["active"].(bool),
			})
		}

		request := &nextdns.CreateDenylistRequest{
			ProfileID: profileID,
			Denylist:  denylist,
		}
		err := client.Denylist.Create(ctx, request)
		if err != nil {
			return diag.FromErr(errors.Wrap(err, "error creating deny list"))
		}
	}

	return resourceNextDNSDenylistRead(ctx, d, meta)
}

func resourceNextDNSDenylistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.CreateDenylistRequest{
		ProfileID: profileID,
		Denylist:  []*nextdns.Denylist{},
	}
	err := client.Denylist.Create(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting deny list"))
	}

	return resourceNextDNSDenylistRead(ctx, d, meta)
}

func resourceNextDNSDenylistImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSDenylistRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
