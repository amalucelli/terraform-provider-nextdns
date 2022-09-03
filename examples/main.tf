resource "nextdns_denylist" "this" {
  profile_id = "3c9e29"

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
  profile_id = "3c9e29"

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
  profile_id = "3c9e29"

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

terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}
