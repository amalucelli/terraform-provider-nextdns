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

func resourceNextDNSProfile() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSProfileSchema(),
		CreateContext: resourceNextDNSProfileCreate,
		ReadContext:   resourceNextDNSProfileRead,
		UpdateContext: resourceNextDNSProfileUpdate,
		DeleteContext: resourceNextDNSProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSProfileImport,
		},
	}
}

func resourceNextDNSProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)

	request := &nextdns.CreateProfileRequest{
		Name: d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	profileID, err := client.Profiles.Create(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating profile"))
	}

	d.SetId(profileID)
	d.Set("profile_id", profileID)

	return resourceNextDNSProfileRead(ctx, d, meta)
}

func resourceNextDNSProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.GetProfileRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	profile, err := client.Profiles.Get(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error getting profile"))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", profile))

	if err := d.Set("name", profile.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(profileID)

	endpoint := map[string]interface{}{}
	endpoint["doh"] = DNSOverHTTPSAddress(profileID)
	endpoint["dot"] = DNSOverTLSAddress(profileID)
	endpoint["ipv4"] = profile.Setup.Ipv4
	endpoint["ipv6"] = profile.Setup.Ipv6

	d.Set("endpoint", []map[string]interface{}{endpoint})

	linkedip := map[string]interface{}{}
	linkedip["servers"] = profile.Setup.LinkedIP.Servers
	linkedip["ip"] = profile.Setup.LinkedIP.IP
	linkedip["update_token"] = profile.Setup.LinkedIP.UpdateToken

	d.Set("linkedip", []map[string]interface{}{linkedip})

	return nil
}

func resourceNextDNSProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	profile := &nextdns.Profile{
		Name: d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", profile))

	request := &nextdns.UpdateProfileRequest{
		ProfileID: profileID,
		Profile:   profile,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err := client.Profiles.Update(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error updating profile"))
	}

	return resourceNextDNSProfileRead(ctx, d, meta)
}

func resourceNextDNSProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.DeleteProfileRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	err := client.Profiles.Delete(ctx, request)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting profile"))
	}

	return nil
}

func resourceNextDNSProfileImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSProfileRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
