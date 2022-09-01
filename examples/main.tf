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

terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}
