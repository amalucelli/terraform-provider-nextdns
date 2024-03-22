---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nextdns_setup_endpoint Data Source - terraform-provider-nextdns"
subcategory: ""
description: |-
  
---

# nextdns_setup_endpoint (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `profile_id` (String) The profile identifier to target the resource.

### Read-Only

- `dnscrypt` (String) The DNS Stamps from the profile.
- `doh` (String) The DNS over HTTPS address the profile is reachable at.
- `dot` (String) The DNS over TLS address the profile is reachable at.
- `id` (String) The ID of this resource.
- `ipv4` (List of String) The IPv4 addresses the profile is reachable at.
- `ipv6` (List of String) The IPv6 addresses the profile is reachable at.