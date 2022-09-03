variable "profile_id" {
  type    = string
  default = "abc123"
}

resource "nextdns_denylist" "this" {
  profile_id = var.profile_id

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
  profile_id = var.profile_id

  domain {
    id     = "search.brave.com"
    active = true
  }

  domain {
    id     = "duckduckgo.com"
    active = false
  }
}

resource "nextdns_security" "this" {
  profile_id = var.profile_id

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

  tlds = [
    "pizza",
    "beer",
    "meme",
    "ninja",
  ]
}

resource "nextdns_privacy" "this" {
  profile_id = var.profile_id

  disguised_trackers = true
  allow_affiliate    = false

  blocklists = [
    "nextdns-recommended",
    "easylist",
    "goodbye-ads",
    "1hosts-pro",
    "easyprivacy",
    "adguard-base-filter",
  ]

  natives = [
    "apple",
    "alexa",
    "windows",
    "xiaomi",
    "huawei",
    "roku",
    "sonos",
    "samsung",
  ]
}

terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}
