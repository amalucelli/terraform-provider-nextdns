package nextdns

import (
	"context"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceNextDNSParentalControl() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSParentalControlSchema(),
		CreateContext: resourceNextDNSParentalControlCreate,
		ReadContext:   resourceNextDNSParentalControlRead,
		UpdateContext: resourceNextDNSParentalControlUpdate,
		DeleteContext: resourceNextDNSParentalControlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSParentalControlImport,
		},
	}
}

func resourceNextDNSParentalControlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	parentalControl, err := buildParentalControl(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error building parental control settings from resource"))
	}

	request := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: parentalControl,
	}
	err = client.ParentalControl.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating parental control settings"))
	}

	d.SetId(profileID)

	return resourceNextDNSParentalControlRead(ctx, d, meta)
}

func resourceNextDNSParentalControlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetParentalControlRequest{
		ProfileID: profileID,
	}
	parentalControl, err := client.ParentalControl.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting parental control settings"))
	}

	var services []map[string]interface{}

	for _, s := range parentalControl.Services {
		service := make(map[string]interface{})
		service["id"] = s.ID
		service["active"] = s.Active

		services = append(services, service)
	}
	if err := d.Set("service", services); err != nil {
		return diag.FromErr(err)
	}

	var categories []map[string]interface{}

	for _, c := range parentalControl.Categories {
		category := make(map[string]interface{})
		category["id"] = c.ID
		category["active"] = c.Active

		categories = append(categories, category)
	}
	if err := d.Set("category", categories); err != nil {
		return diag.FromErr(err)
	}

	d.Set("block_bypass", parentalControl.BlockBypass)
	d.Set("safe_search", parentalControl.SafeSearch)
	d.Set("youtube_restricted_mode", parentalControl.YoutubeRestrictedMode)

	d.SetId(profileID)

	return nil
}

func resourceNextDNSParentalControlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	parentalControl, err := buildParentalControl(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating parental control settings"))
	}

	request := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: parentalControl,
	}
	err = client.ParentalControl.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating parental control settings"))
	}

	return resourceNextDNSParentalControlRead(ctx, d, meta)
}

func resourceNextDNSParentalControlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	services := &nextdns.CreateParentalControlServicesRequest{
		ProfileID:               profileID,
		ParentalControlServices: []*nextdns.ParentalControlServices{},
	}
	err := client.ParentalControlServices.Create(ctx, services)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting services settings"))
	}

	categories := &nextdns.CreateParentalControlCategoriesRequest{
		ProfileID:                 profileID,
		ParentalControlCategories: []*nextdns.ParentalControlCategories{},
	}
	err = client.ParentalControlCategories.Create(ctx, categories)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting categories settings"))
	}

	parentalControl := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: &nextdns.ParentalControl{},
	}
	err = client.ParentalControl.Update(ctx, parentalControl)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting parental control settings"))
	}

	return resourceNextDNSParentalControlRead(ctx, d, meta)

}

func resourceNextDNSParentalControlImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSParentalControlRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func buildParentalControl(d *schema.ResourceData) (*nextdns.ParentalControl, error) {
	ParentalControl := &nextdns.ParentalControl{
		BlockBypass:           d.Get("block_bypass").(bool),
		SafeSearch:            d.Get("safe_search").(bool),
		YoutubeRestrictedMode: d.Get("youtube_restricted_mode").(bool),
	}

	foundSvc, ok := d.GetOk("service")
	if !ok {
		return nil, errors.New("unable to find service in resource data")
	}

	recordsSvc := foundSvc.(*schema.Set).List()

	services := make([]*nextdns.ParentalControlServices, len(recordsSvc))
	for k, v := range recordsSvc {
		services[k] = &nextdns.ParentalControlServices{
			ID:     v.(map[string]interface{})["id"].(string),
			Active: v.(map[string]interface{})["active"].(bool),
		}
	}
	ParentalControl.Services = services

	foundCat, ok := d.GetOk("category")
	if !ok {
		return nil, errors.New("unable to find category in resource data")
	}

	recordsCat := foundCat.(*schema.Set).List()

	categories := make([]*nextdns.ParentalControlCategories, len(recordsCat))
	for k, v := range recordsCat {
		categories[k] = &nextdns.ParentalControlCategories{
			ID:     v.(map[string]interface{})["id"].(string),
			Active: v.(map[string]interface{})["active"].(bool),
		}
	}
	ParentalControl.Categories = categories

	return ParentalControl, nil
}
