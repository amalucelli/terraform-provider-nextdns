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

resource "nextdns_parental_control" "this" {
  profile_id = var.profile_id

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

resource "nextdns_settings" "this" {
  profile_id = var.profile_id

  logs {
    enabled = true

    privacy {
      log_clients_ip = true
      log_domains    = true
    }

    retention = 7776000
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

terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}
