package nextdns

import (
	"context"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		return diag.FromErr(fmt.Errorf("error creating parental control settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", parentalControl))

	services := &nextdns.CreateParentalControlServicesRequest{
		ProfileID:               profileID,
		ParentalControlServices: parentalControl.Services,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", services))

	err = client.ParentalControlServices.Create(ctx, services)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating services settings: %w", err))
	}

	categories := &nextdns.CreateParentalControlCategoriesRequest{
		ProfileID:                 profileID,
		ParentalControlCategories: parentalControl.Categories,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", categories))

	err = client.ParentalControlCategories.Create(ctx, categories)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating categories settings: %w", err))
	}

	request := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: parentalControl,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.ParentalControl.Update(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating parental control settings: %w", err))
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
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	parentalControl, err := client.ParentalControl.Get(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting parental control settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", parentalControl))

	if parentalControl.Recreation != nil {
		recreation := map[string]interface{}{}
		recreation["timezone"] = parentalControl.Recreation.Timezone

		if parentalControl.Recreation.Times.Monday != nil {
			recreation["monday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Monday.Start,
					"end":   parentalControl.Recreation.Times.Monday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Tuesday != nil {
			recreation["tuesday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Tuesday.Start,
					"end":   parentalControl.Recreation.Times.Tuesday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Wednesday != nil {
			recreation["wednesday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Wednesday.Start,
					"end":   parentalControl.Recreation.Times.Wednesday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Thursday != nil {
			recreation["thursday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Thursday.Start,
					"end":   parentalControl.Recreation.Times.Thursday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Friday != nil {
			recreation["friday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Friday.Start,
					"end":   parentalControl.Recreation.Times.Friday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Saturday != nil {
			recreation["saturday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Saturday.Start,
					"end":   parentalControl.Recreation.Times.Saturday.End,
				},
			}
		}
		if parentalControl.Recreation.Times.Sunday != nil {
			recreation["sunday"] = []map[string]interface{}{
				{
					"start": parentalControl.Recreation.Times.Sunday.Start,
					"end":   parentalControl.Recreation.Times.Sunday.End,
				},
			}
		}

		d.Set("recreation", []map[string]interface{}{recreation})
	}

	var services []map[string]interface{}
	for _, s := range parentalControl.Services {
		service := make(map[string]interface{})
		service["id"] = s.ID
		service["active"] = s.Active
		service["recreation"] = s.Recreation

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
		category["recreation"] = c.Recreation

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
		return diag.FromErr(fmt.Errorf("error updating parental control settings: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", parentalControl))

	services := &nextdns.CreateParentalControlServicesRequest{
		ProfileID:               profileID,
		ParentalControlServices: parentalControl.Services,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", services))

	err = client.ParentalControlServices.Create(ctx, services)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating services settings: %w", err))
	}

	categories := &nextdns.CreateParentalControlCategoriesRequest{
		ProfileID:                 profileID,
		ParentalControlCategories: parentalControl.Categories,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", categories))

	err = client.ParentalControlCategories.Create(ctx, categories)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating categories settings: %w", err))
	}

	request := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: parentalControl,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.ParentalControl.Update(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating parental control settings: %w", err))
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
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", services))

	err := client.ParentalControlServices.Create(ctx, services)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting services settings: %w", err))
	}

	categories := &nextdns.CreateParentalControlCategoriesRequest{
		ProfileID:                 profileID,
		ParentalControlCategories: []*nextdns.ParentalControlCategories{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", categories))

	err = client.ParentalControlCategories.Create(ctx, categories)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting categories settings: %w", err))
	}

	parentalControl := &nextdns.UpdateParentalControlRequest{
		ProfileID:       profileID,
		ParentalControl: &nextdns.ParentalControl{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", parentalControl))

	err = client.ParentalControl.Update(ctx, parentalControl)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting parental control settings: %w", err))
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

	ParentalControl.Recreation = &nextdns.ParentalControlRecreation{}
	if _, ok := d.GetOk("recreation"); ok {
		times := &nextdns.ParentalControlRecreationTimes{}

		if _, ok := d.GetOk("recreation.0.monday"); ok {
			times.Monday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.monday.0.start").(string),
				End:   d.Get("recreation.0.monday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.tuesday"); ok {
			times.Tuesday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.tuesday.0.start").(string),
				End:   d.Get("recreation.0.tuesday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.wednesday"); ok {
			times.Wednesday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.wednesday.0.start").(string),
				End:   d.Get("recreation.0.wednesday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.thursday"); ok {
			times.Thursday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.thursday.0.start").(string),
				End:   d.Get("recreation.0.thursday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.friday"); ok {
			times.Friday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.friday.0.start").(string),
				End:   d.Get("recreation.0.friday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.saturday"); ok {
			times.Saturday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.saturday.0.start").(string),
				End:   d.Get("recreation.0.saturday.0.end").(string),
			}
		}

		if _, ok := d.GetOk("recreation.0.sunday"); ok {
			times.Sunday = &nextdns.ParentalControlRecreationInterval{
				Start: d.Get("recreation.0.sunday.0.start").(string),
				End:   d.Get("recreation.0.sunday.0.end").(string),
			}
		}

		recreation := &nextdns.ParentalControlRecreation{
			Times:    times,
			Timezone: d.Get("recreation.0.timezone").(string),
		}

		ParentalControl.Recreation = recreation
	}

	ParentalControl.Services = []*nextdns.ParentalControlServices{}
	if foundSvc, ok := d.GetOk("service"); ok {
		recordsSvc := foundSvc.(*schema.Set).List()
		services := make([]*nextdns.ParentalControlServices, len(recordsSvc))

		for k, v := range recordsSvc {
			services[k] = &nextdns.ParentalControlServices{
				ID:         v.(map[string]interface{})["id"].(string),
				Active:     v.(map[string]interface{})["active"].(bool),
				Recreation: v.(map[string]interface{})["recreation"].(bool),
			}
		}
		ParentalControl.Services = services
	}

	ParentalControl.Categories = []*nextdns.ParentalControlCategories{}
	if foundCat, ok := d.GetOk("category"); ok {
		recordsCat := foundCat.(*schema.Set).List()
		categories := make([]*nextdns.ParentalControlCategories, len(recordsCat))

		for k, v := range recordsCat {
			categories[k] = &nextdns.ParentalControlCategories{
				ID:         v.(map[string]interface{})["id"].(string),
				Active:     v.(map[string]interface{})["active"].(bool),
				Recreation: v.(map[string]interface{})["recreation"].(bool),
			}
		}
		ParentalControl.Categories = categories
	}

	return ParentalControl, nil
}
