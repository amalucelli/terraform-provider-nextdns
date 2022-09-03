# terraform-provider-nextdns

[NextDNS](https://nextdns.io/) provider for [Terraform](https://terraform.io).

This provider lets you declaratively define the configuration for your NextDNS profiles.

## Requirements

An API Key is required to interact with the NextDNS API.
You can find your API Key in the [NextDNS account](https://my.nextdns.io/account) page.

## Getting Started

**Important**: This provider is not yet published on the Terraform Registry
and it is still under development. If you want to use it, you will have to build and install it locally.

```hcl
terraform {
  required_providers {
    nextdns = {
      source  = "github.com/amalucelli/nextdns"
      version = "0.1.0"
    }
  }
}

provider "nextdns" {
  api_key = "NEXTDNS_API_KEY"
}
```
