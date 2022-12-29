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

	allowlist, err := buildAllowlist(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error building allow list"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", allowlist))

	request := &nextdns.CreateAllowlistRequest{
		ProfileID: profileID,
		Allowlist: allowlist,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Allowlist.Create(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating allow list"))
	}

	d.SetId(profileID)

	return resourceNextDNSAllowlistRead(ctx, d, meta)
}

func resourceNextDNSAllowlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.ListAllowlistRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	allowlist, err := client.Allowlist.List(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting allow list"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", allowlist))

	var domains []map[string]interface{}

	for _, d := range allowlist {
		domain := make(map[string]interface{})
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

	allowlist, err := buildAllowlist(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error building allow list"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", allowlist))

	request := &nextdns.CreateAllowlistRequest{
		ProfileID: profileID,
		Allowlist: allowlist,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Allowlist.Create(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating allow list"))
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
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

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

func buildAllowlist(d *schema.ResourceData) ([]*nextdns.Allowlist, error) {
	found, ok := d.GetOk("domain")
	if !ok {
		return nil, errors.New("unable to find domain in resource data")
	}

	records := found.(*schema.Set).List()

	allowlist := make([]*nextdns.Allowlist, len(records))
	for k, v := range records {
		allowlist[k] = &nextdns.Allowlist{
			ID:     v.(map[string]interface{})["id"].(string),
			Active: v.(map[string]interface{})["active"].(bool),
		}
	}

	return allowlist, nil
}
