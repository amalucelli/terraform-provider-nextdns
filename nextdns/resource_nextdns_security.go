package nextdns

import (
	"context"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSSecurity() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSSecuritySchema(),
		CreateContext: resourceNextDNSSecurityCreate,
		ReadContext:   resourceNextDNSSecurityRead,
		UpdateContext: resourceNextDNSSecurityUpdate,
		DeleteContext: resourceNextDNSSecurityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSSecurityImport,
		},
	}
}

func resourceNextDNSSecurityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	sec, err := buildSecurity(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating security settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", sec))

	tlds := &nextdns.CreateSecurityTldsRequest{
		ProfileID:    profileID,
		SecurityTlds: sec.Tlds,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", tlds))

	err = client.SecurityTlds.Create(ctx, tlds)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating security tlds settings: %w", err))
	}

	request := &nextdns.UpdateSecurityRequest{
		ProfileID: profileID,
		Security:  sec,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Security.Update(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating security settings: %w", err))
	}

	d.SetId(profileID)

	return resourceNextDNSSecurityRead(ctx, d, meta)
}

func resourceNextDNSSecurityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetSecurityRequest{
		ProfileID: profileID,
	}
	security, err := client.Security.Get(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting security settings: %w", err))
	}

	d.SetId(profileID)

	d.Set("threat_intelligence_feeds", security.ThreatIntelligenceFeeds)
	d.Set("ai_threat_detection", security.AiThreatDetection)
	d.Set("google_safe_browsing", security.GoogleSafeBrowsing)
	d.Set("crypto_jacking", security.Cryptojacking)
	d.Set("dns_rebinding", security.DNSRebinding)
	d.Set("idn_homographs", security.IdnHomographs)
	d.Set("typo_squatting", security.Typosquatting)
	d.Set("dga", security.Dga)
	d.Set("nrd", security.Nrd)
	d.Set("ddns", security.DDNS)
	d.Set("parking", security.Parking)
	d.Set("csam", security.Csam)

	d.Set("tlds", flattenTLDs(security.Tlds))
	return nil
}

func resourceNextDNSSecurityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	sec, err := buildSecurity(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating security settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", sec))

	tlds := &nextdns.CreateSecurityTldsRequest{
		ProfileID:    profileID,
		SecurityTlds: sec.Tlds,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", tlds))

	err = client.SecurityTlds.Create(ctx, tlds)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating security tlds settings: %w", err))
	}

	request := &nextdns.UpdateSecurityRequest{
		ProfileID: profileID,
		Security:  sec,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Security.Update(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating security settings: %w", err))
	}

	return resourceNextDNSSecurityRead(ctx, d, meta)
}

func resourceNextDNSSecurityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	tlds := &nextdns.CreateSecurityTldsRequest{
		ProfileID:    profileID,
		SecurityTlds: []*nextdns.SecurityTlds{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", tlds))

	err := client.SecurityTlds.Create(ctx, tlds)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting security tlds settings: %w", err))
	}

	sec := &nextdns.UpdateSecurityRequest{
		ProfileID: profileID,
		Security:  &nextdns.Security{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", sec))

	err = client.Security.Update(ctx, sec)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting security settings: %w", err))
	}

	return resourceNextDNSSecurityRead(ctx, d, meta)
}

func resourceNextDNSSecurityImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSSecurityRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func flattenTLDs(tlds []*nextdns.SecurityTlds) []string {
	ids := make([]string, 0)
	for _, tld := range tlds {
		ids = append(ids, tld.ID)
	}

	return ids
}

func buildSecurity(d *schema.ResourceData) (*nextdns.Security, error) {
	sec := &nextdns.Security{
		ThreatIntelligenceFeeds: d.Get("threat_intelligence_feeds").(bool),
		AiThreatDetection:       d.Get("ai_threat_detection").(bool),
		GoogleSafeBrowsing:      d.Get("google_safe_browsing").(bool),
		Cryptojacking:           d.Get("crypto_jacking").(bool),
		DNSRebinding:            d.Get("dns_rebinding").(bool),
		IdnHomographs:           d.Get("idn_homographs").(bool),
		Typosquatting:           d.Get("typo_squatting").(bool),
		Dga:                     d.Get("dga").(bool),
		Nrd:                     d.Get("nrd").(bool),
		DDNS:                    d.Get("ddns").(bool),
		Parking:                 d.Get("parking").(bool),
		Csam:                    d.Get("csam").(bool),
	}

	sec.Tlds = []*nextdns.SecurityTlds{}
	if found, ok := d.GetOk("tlds"); ok {
		recordsTlds := found.([]interface{})
		tlds := make([]*nextdns.SecurityTlds, len(recordsTlds))

		for k, v := range recordsTlds {
			tlds[k] = &nextdns.SecurityTlds{
				ID: v.(string),
			}
		}
		sec.Tlds = tlds
	}

	return sec, nil
}
