resource "nextdns_profile" "this" {
  name = "terraform-provider-nextdns"
}

resource "nextdns_denylist" "this" {
  profile_id = nextdns_profile.this.id

  domain {
    id     = "google.com"
    active = true
  }

  domain {
    id     = "bing.com"
    active = false
  }
}

resource "nextdns_allowlist" "this" {
  profile_id = nextdns_profile.this.id

  domain {
    id     = "search.brave.com"
    active = true
  }

  domain {
    id     = "duckduckgo.com"
    active = false
  }
}

resource "nextdns_parental_control" "this" {
  profile_id = nextdns_profile.this.id

  safe_search             = true
  youtube_restricted_mode = false
  block_bypass            = true

  service {
    id         = "tiktok"
    active     = true
    recreation = false
  }

  service {
    id         = "instagram"
    active     = false
    recreation = false
  }

  service {
    id         = "facebook"
    active     = true
    recreation = true
  }

  category {
    id         = "dating"
    active     = false
    recreation = true
  }

  category {
    id         = "gambling"
    active     = false
    recreation = false
  }

  recreation {
    timezone = "America/New_York"

    monday {
      start = "07:02:00"
      end   = "23:00:00"
    }
    saturday {
      start = "20:22:00"
      end   = "23:33:00"
    }
    friday {
      start = "12:00:00"
      end   = "15:35:00"
    }
    sunday {
      start = "01:00:00"
      end   = "23:00:00"
    }
  }
}

resource "nextdns_security" "this" {
  profile_id = nextdns_profile.this.id

  threat_intelligence_feeds = true
  ai_threat_detection       = false
  google_safe_browsing      = true
  crypto_jacking            = false
  dns_rebinding             = true
  idn_homographs            = false
  typo_squatting            = true
  dga                       = false
  nrd                       = true
  ddns                      = false
  parking                   = true
  csam                      = false

  tlds = toset([
    "pizza",
    "beer",
    "meme",
    "ninja",
  ])
}

resource "nextdns_privacy" "this" {
  profile_id = nextdns_profile.this.id

  disguised_trackers = true
  allow_affiliate    = false

  blocklists = toset([
    "nextdns-recommended",
    "easylist",
    "goodbye-ads",
    "1hosts-pro",
    "easyprivacy",
    "adguard-base-filter",
  ])

  natives = toset([
    "apple",
    "alexa",
    "windows",
    "xiaomi",
    "huawei",
    "roku",
    "sonos",
    "samsung",
  ])
}

resource "nextdns_settings" "this" {
  profile_id = nextdns_profile.this.id

  logs {
    enabled = true

    privacy {
      log_clients_ip = true
      log_domains    = true
    }

    retention = "1 day"
    location  = "us"
  }

  block_page {
    enabled = true
  }

  performance {
    ecs              = true
    cache_boost      = true
    cname_flattening = true
  }

  web3 = true
}

resource "nextdns_rewrite" "this" {
  profile_id = nextdns_profile.this.id

  rewrite {
    domain  = "google.com"
    address = "1.1.1.1"
  }

  rewrite {
    domain  = "cloudflare.com"
    address = "8.8.8.8"
  }
}

data "nextdns_setup_endpoint" "this" {
  profile_id = nextdns_profile.this.id
}

data "nextdns_setup_linkedip" "this" {
  profile_id = nextdns_profile.this.id
}

terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}

output "doh" {
  description = "The DNS over HTTPS address the profile is reachable at"
  value = data.nextdns_setup_endpoint.this.dot
}

output "dot" {
  description = "The DNS over TLS address the profile is reachable at"
  value = data.nextdns_setup_endpoint.this.doh
}

output "ipv6" {
  description = "The IPv6 address the profile is reachable at"
  value = data.nextdns_setup_endpoint.this.ipv6
}

output "servers" {
  description = "The DNS servers available for the profile"
  value = data.nextdns_setup_linkedip.this.servers
}
