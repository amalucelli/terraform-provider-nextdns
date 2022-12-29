package nextdns

import (
	"context"
	"fmt"
	"time"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceNextDNSSettings() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSSettingsSchema(),
		CreateContext: resourceNextDNSSettingsCreate,
		ReadContext:   resourceNextDNSSettingsRead,
		UpdateContext: resourceNextDNSSettingsUpdate,
		DeleteContext: resourceNextDNSSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSSettingsImport,
		},
	}
}

func resourceNextDNSSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	settings, err := buildSettings(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating settings"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", settings))

	logs := &nextdns.UpdateSettingsLogsRequest{
		ProfileID:    profileID,
		SettingsLogs: settings.Logs,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", logs))

	err = client.SettingsLogs.Update(ctx, logs)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating logs settings"))
	}

	blockPage := &nextdns.UpdateSettingsBlockPageRequest{
		ProfileID:         profileID,
		SettingsBlockPage: settings.BlockPage,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", blockPage))

	err = client.SettingsBlockPage.Update(ctx, blockPage)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating categories settings"))
	}

	performance := &nextdns.UpdateSettingsPerformanceRequest{
		ProfileID:           profileID,
		SettingsPerformance: settings.Performance,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", performance))

	err = client.SettingsPerformance.Update(ctx, performance)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating categories settings"))
	}

	request := &nextdns.UpdateSettingsRequest{
		ProfileID: profileID,
		Settings:  settings,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Settings.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating parental control settings"))
	}

	d.SetId(profileID)

	return resourceNextDNSSettingsRead(ctx, d, meta)
}

func resourceNextDNSSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetSettingsRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	settings, err := client.Settings.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting settings"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", settings))

	logs := map[string]interface{}{}
	logs["enabled"] = settings.Logs.Enabled
	logs["privacy"] = []map[string]interface{}{
		{
			"log_clients_ip": invertPrivacySettings(settings.Logs.Drop.IP),
			"log_domains":    invertPrivacySettings(settings.Logs.Drop.Domain),
		},
	}
	logs["retention"] = convertSecondsToRetention(settings.Logs.Retention)
	logs["location"] = settings.Logs.Location

	d.Set("logs", []map[string]interface{}{logs})

	blockPage := map[string]interface{}{}
	blockPage["enabled"] = settings.BlockPage.Enabled

	d.Set("block_page", []map[string]interface{}{blockPage})

	performance := map[string]interface{}{}
	performance["ecs"] = settings.Performance.Ecs
	performance["cache_boost"] = settings.Performance.CacheBoost
	performance["cname_flattening"] = settings.Performance.CnameFlattening

	d.Set("performance", []map[string]interface{}{performance})

	d.Set("web3", settings.Web3)

	d.SetId(profileID)

	return nil
}

func resourceNextDNSSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	settings, err := buildSettings(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating settings"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", settings))

	logs := &nextdns.UpdateSettingsLogsRequest{
		ProfileID:    profileID,
		SettingsLogs: settings.Logs,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", logs))

	err = client.SettingsLogs.Update(ctx, logs)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating logs settings"))
	}

	blockPage := &nextdns.UpdateSettingsBlockPageRequest{
		ProfileID:         profileID,
		SettingsBlockPage: settings.BlockPage,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", blockPage))

	err = client.SettingsBlockPage.Update(ctx, blockPage)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating categories settings"))
	}

	performance := &nextdns.UpdateSettingsPerformanceRequest{
		ProfileID:           profileID,
		SettingsPerformance: settings.Performance,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", performance))

	err = client.SettingsPerformance.Update(ctx, performance)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating categories settings"))
	}

	request := &nextdns.UpdateSettingsRequest{
		ProfileID: profileID,
		Settings:  settings,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err = client.Settings.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating parental control settings"))
	}

	return resourceNextDNSSettingsRead(ctx, d, meta)
}
func resourceNextDNSSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	logs := &nextdns.UpdateSettingsLogsRequest{
		ProfileID:    profileID,
		SettingsLogs: &nextdns.SettingsLogs{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", logs))

	err := client.SettingsLogs.Update(ctx, logs)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting logs settings"))
	}

	blockPage := &nextdns.UpdateSettingsBlockPageRequest{
		ProfileID:         profileID,
		SettingsBlockPage: &nextdns.SettingsBlockPage{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", blockPage))

	err = client.SettingsBlockPage.Update(ctx, blockPage)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting block page settings"))
	}

	performance := &nextdns.UpdateSettingsPerformanceRequest{
		ProfileID:           profileID,
		SettingsPerformance: &nextdns.SettingsPerformance{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", performance))

	err = client.SettingsPerformance.Update(ctx, performance)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting performance settings"))
	}

	settings := &nextdns.UpdateSettingsRequest{
		ProfileID: profileID,
		Settings:  &nextdns.Settings{},
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", settings))

	err = client.Settings.Update(ctx, settings)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting settings"))
	}

	return resourceNextDNSSettingsRead(ctx, d, meta)
}

func resourceNextDNSSettingsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSSettingsRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func buildSettings(d *schema.ResourceData) (*nextdns.Settings, error) {
	logs := &nextdns.SettingsLogs{
		Enabled: d.Get("logs.0.enabled").(bool),
		Drop: &nextdns.SettingsLogsDrop{
			IP:     invertPrivacySettings(d.Get("logs.0.privacy.0.log_clients_ip").(bool)),
			Domain: invertPrivacySettings(d.Get("logs.0.privacy.0.log_domains").(bool)),
		},
		Retention: convertRetentionToSeconds(d.Get("logs.0.retention").(string)),
		Location:  d.Get("logs.0.location").(string),
	}

	blockPage := &nextdns.SettingsBlockPage{
		Enabled: d.Get("block_page.0.enabled").(bool),
	}

	performance := &nextdns.SettingsPerformance{
		Ecs:             d.Get("performance.0.ecs").(bool),
		CacheBoost:      d.Get("performance.0.cache_boost").(bool),
		CnameFlattening: d.Get("performance.0.cname_flattening").(bool),
	}

	Settings := &nextdns.Settings{
		Logs:        logs,
		BlockPage:   blockPage,
		Performance: performance,
		Web3:        d.Get("web3").(bool),
	}

	return Settings, nil
}

func convertRetentionToSeconds(retention string) int {
	switch retention {
	case "1 hour":
		d := time.Hour
		return int(d.Seconds())
	case "6 hours":
		d := time.Hour * 6
		return int(d.Seconds())
	case "1 day":
		d := time.Hour * 24
		return int(d.Seconds())
	case "1 week":
		d := time.Hour * 24 * 7
		return int(d.Seconds())
	case "1 month":
		d := time.Hour * 24 * 30
		return int(d.Seconds())
	case "3 months":
		d := time.Hour * 24 * 90
		return int(d.Seconds())
	case "6 months":
		d := time.Hour * 24 * 180
		return int(d.Seconds())
	case "1 year":
		d := time.Hour * 24 * 365
		return int(d.Seconds())
	case "2 years":
		d := time.Hour * 24 * 365 * 2
		return int(d.Seconds())
	default:
		return 0
	}
}

func convertSecondsToRetention(seconds int) string {
	switch seconds {
	case int(time.Hour.Seconds()):
		return "1 hour"
	case int(time.Hour.Seconds() * 6):
		return "6 hours"
	case int(time.Hour.Seconds() * 24):
		return "1 day"
	case int(time.Hour.Seconds() * 24 * 7):
		return "1 week"
	case int(time.Hour.Seconds() * 24 * 30):
		return "1 month"
	case int(time.Hour.Seconds() * 24 * 90):
		return "3 months"
	case int(time.Hour.Seconds() * 24 * 180):
		return "6 months"
	case int(time.Hour.Seconds() * 24 * 365):
		return "1 year"
	case int(time.Hour.Seconds() * 24 * 365 * 2):
		return "2 years"
	default:
		return ""
	}
}

// invertPrivacySettings inverts the privacy settings,
// as the API wants the opposite of what the user wants.
func invertPrivacySettings(value bool) bool {
	return !value
}
