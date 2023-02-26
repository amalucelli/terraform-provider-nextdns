# Examples

This directory contains examples of how to use the `terraform-provider-nextdns` provider.

## Usage

To run this example, first create a `providers.tf` file with the following content:

```hcl
provider "nextdns" {
  api_key = "API_KEY"
}
```

Then, apply it with:

```bash
terraform init
terraform apply
```

This will create a new profile called `terraform-provider-nextdns` with the configurations defined in `main.tf`.

You can destroy the profile created with:

```bash
terraform destroy
```
