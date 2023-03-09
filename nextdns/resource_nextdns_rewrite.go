package nextdns

import (
	"context"
	"errors"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNextDNSRewrite() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceNextDNSRewriteSchema(),
		CreateContext: resourceNextDNSRewriteCreate,
		ReadContext:   resourceNextDNSRewriteRead,
		UpdateContext: resourceNextDNSRewriteUpdate,
		DeleteContext: resourceNextDNSRewriteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNextDNSRewriteImport,
		},
	}
}

func resourceNextDNSRewriteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	rewrites, err := buildRewrite(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error building rewrite list: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	list := &nextdns.ListRewritesRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", list))

	existing, err := client.Rewrites.List(ctx, list)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting rewrites: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	// Since an object can not be created twice, we need to delete it first the ones that already exist and are declared in the terraform file.
	// The ones that only exist in the API will still be there.
	for _, e := range existing {
		for _, r := range rewrites {
			if e.Name == r.Name && e.Content == r.Content {
				tflog.Debug(ctx, fmt.Sprintf("rewrite already exists: %+v", e))

				deleteRequest := &nextdns.DeleteRewritesRequest{
					ProfileID: profileID,
					ID:        e.ID,
				}
				tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", deleteRequest))

				err := client.Rewrites.Delete(ctx, deleteRequest)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error deleting rewrite: %w", err))
				}

				continue
			}
		}
	}

	for _, rewrite := range rewrites {
		request := &nextdns.CreateRewritesRequest{
			ProfileID: profileID,
			Rewrites:  rewrite,
		}
		tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

		_, err = client.Rewrites.Create(ctx, request)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating rewrite: %w", err))
		}
	}

	d.SetId(profileID)

	return resourceNextDNSRewriteRead(ctx, d, meta)
}

func resourceNextDNSRewriteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.ListRewritesRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	rewrites, err := client.Rewrites.List(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting rewrites: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	var rewrite []map[string]interface{}

	for _, d := range rewrites {
		record := make(map[string]interface{})
		record["domain"] = d.Name
		record["address"] = d.Content

		rewrite = append(rewrite, record)
	}
	if err := d.Set("rewrite", rewrite); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(profileID)

	return nil
}

func resourceNextDNSRewriteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	rewrites, err := buildRewrite(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error building rewrite list: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	list := &nextdns.ListRewritesRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", list))

	existing, err := client.Rewrites.List(ctx, list)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting rewrites: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	var toRemove []*nextdns.Rewrites
	var toAdd []*nextdns.Rewrites

	for _, e := range existing {
		found := false
		for _, r := range rewrites {
			if e.Name == r.Name && e.Content == r.Content {
				found = true
				tflog.Debug(ctx, fmt.Sprintf("rewrite already exists: %+v", e))
			}
		}
		if !found {
			toRemove = append(toRemove, e)
		}
	}

	for _, r := range rewrites {
		found := false
		for _, e := range existing {
			if e.Name == r.Name && e.Content == r.Content {
				found = true
				tflog.Debug(ctx, fmt.Sprintf("rewrite already exists: %+v", e))
			}
		}
		if !found {
			toAdd = append(toAdd, r)
		}
	}

	for _, r := range toRemove {
		deleteRequest := &nextdns.DeleteRewritesRequest{
			ProfileID: profileID,
			ID:        r.ID,
		}
		tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", deleteRequest))

		err := client.Rewrites.Delete(ctx, deleteRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting rewrite: %w", err))
		}
	}

	for _, r := range toAdd {
		request := &nextdns.CreateRewritesRequest{
			ProfileID: profileID,
			Rewrites:  r,
		}
		tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

		_, err = client.Rewrites.Create(ctx, request)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating rewrite: %w", err))
		}
	}

	return resourceNextDNSRewriteRead(ctx, d, meta)
}

func resourceNextDNSRewriteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nextdns.Client)
	profileID := d.Get("profile_id").(string)

	request := &nextdns.ListRewritesRequest{
		ProfileID: profileID,
	}
	tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

	rewrites, err := client.Rewrites.List(ctx, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting rewrites: %w", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("object built: %+v", rewrites))

	for _, rewrite := range rewrites {
		request := &nextdns.DeleteRewritesRequest{
			ProfileID: profileID,
			ID:        rewrite.ID,
		}
		tflog.Debug(ctx, fmt.Sprintf("request to nextdns api: %+v", request))

		err = client.Rewrites.Delete(ctx, request)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting rewrite: %w", err))
		}
	}

	return resourceNextDNSRewriteRead(ctx, d, meta)
}

func resourceNextDNSRewriteImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	d.SetId(profileID)
	d.Set("profile_id", profileID)

	resourceNextDNSRewriteRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func buildRewrite(d *schema.ResourceData) ([]*nextdns.Rewrites, error) {
	found, ok := d.GetOk("rewrite")
	if !ok {
		// nolint:goerr113
		return nil, errors.New("unable to find rewrite in resource data")
	}

	records := found.(*schema.Set).List()

	rewrites := make([]*nextdns.Rewrites, len(records))
	for k, v := range records {
		rewrites[k] = &nextdns.Rewrites{
			Name:    v.(map[string]interface{})["domain"].(string),
			Content: v.(map[string]interface{})["address"].(string),
		}
	}

	return rewrites, nil
}
