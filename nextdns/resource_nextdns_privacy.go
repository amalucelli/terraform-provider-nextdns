package nextdns

import (
	"context"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceNextDNSPrivacy() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSPrivacySchema(),
		CreateContext: resourceNextDNSPrivacyCreate,
		ReadContext:   resourceNextDNSPrivacyRead,
		UpdateContext: resourceNextDNSPrivacyUpdate,
		DeleteContext: resourceNextDNSPrivacyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSPrivacyImport,
		},
	}
}

func resourceNextDNSPrivacyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	privacy, err := buildPrivacy(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error building privacy settings from resource"))
	}

	request := &nextdns.UpdatePrivacyRequest{
		ProfileID: profileID,
		Privacy:   privacy,
	}
	err = client.Privacy.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating privacy settings"))
	}

	d.SetId(profileID)

	return resourceNextDNSPrivacyRead(ctx, d, meta)
}

func resourceNextDNSPrivacyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetPrivacyRequest{
		ProfileID: profileID,
	}
	privacy, err := client.Privacy.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting privacy settings"))
	}

	d.SetId(profileID)

	d.Set("allow_affiliate", privacy.AllowAffiliate)
	d.Set("disguised_trackers", privacy.DisguisedTrackers)
	d.Set("blocklists", flattenBlocklists(privacy.Blocklists))
	d.Set("natives", flattenNatives(privacy.Natives))

	return nil
}

func resourceNextDNSPrivacyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	privacy, err := buildPrivacy(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating privacy settings"))
	}

	request := &nextdns.UpdatePrivacyRequest{
		ProfileID: profileID,
		Privacy:   privacy,
	}
	err = client.Privacy.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating privacy settings"))
	}

	return resourceNextDNSPrivacyRead(ctx, d, meta)
}

func resourceNextDNSPrivacyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	blocklist := &nextdns.CreatePrivacyBlocklistsRequest{
		ProfileID:         profileID,
		PrivacyBlocklists: []*nextdns.PrivacyBlocklists{},
	}
	err := client.PrivacyBlocklists.Create(ctx, blocklist)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting blocklist settings"))
	}

	natives := &nextdns.CreatePrivacyNativesRequest{
		ProfileID:      profileID,
		PrivacyNatives: []*nextdns.PrivacyNatives{},
	}
	err = client.PrivacyNatives.Create(ctx, natives)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting native settings"))
	}

	privacy := &nextdns.UpdatePrivacyRequest{
		ProfileID: profileID,
		Privacy:   &nextdns.Privacy{},
	}
	err = client.Privacy.Update(ctx, privacy)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting privacy settings"))
	}

	return resourceNextDNSPrivacyRead(ctx, d, meta)
}

func resourceNextDNSPrivacyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSPrivacyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func flattenBlocklists(blocklists []*nextdns.PrivacyBlocklists) []string {
	ids := make([]string, 0)
	for _, entry := range blocklists {
		ids = append(ids, entry.ID)
	}

	return ids
}

func flattenNatives(natives []*nextdns.PrivacyNatives) []string {
	ids := make([]string, 0)
	for _, entry := range natives {
		ids = append(ids, entry.ID)
	}

	return ids
}

func buildPrivacy(d *schema.ResourceData) (*nextdns.Privacy, error) {
	privacy := &nextdns.Privacy{
		AllowAffiliate:    d.Get("allow_affiliate").(bool),
		DisguisedTrackers: d.Get("disguised_trackers").(bool),
	}

	rblocklist, ok := d.Get("blocklists").([]interface{})
	if !ok {
		return nil, errors.New("unable to create interface array type assertion")
	}

	blocklist := make([]*nextdns.PrivacyBlocklists, (len(rblocklist)))
	for k, v := range rblocklist {
		blocklist[k] = &nextdns.PrivacyBlocklists{
			ID: v.(string),
		}
	}
	privacy.Blocklists = blocklist

	rnatives, ok := d.Get("natives").([]interface{})
	if !ok {
		return nil, errors.New("unable to create interface array type assertion")
	}

	natives := make([]*nextdns.PrivacyNatives, (len(rnatives)))
	for k, v := range rnatives {
		natives[k] = &nextdns.PrivacyNatives{
			ID: v.(string),
		}
	}
	privacy.Natives = natives

	return privacy, nil
}
