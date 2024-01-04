# terraform-provider-nextdns

[NextDNS](https://nextdns.io/) provider for [Terraform](https://terraform.io).

This provider lets you declaratively define the configuration for your NextDNS profiles.

## Requirements

An API Key is required to interact with the NextDNS API.
You can find your API Key in the [NextDNS account](https://my.nextdns.io/account) page.

## Getting Started

```hcl
terraform {
  required_providers {
    nextdns = {
      source  = "amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}

provider "nextdns" {
  api_key = "NEXTDNS_API_KEY"
}
```
