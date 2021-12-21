---
layout: "twcc"
apage_title: "Provider: TWCC"
description: |-
  The Taiwan Computing Cloud (TWCC) Provider is used to interact with the many resources supported by TWCC.
---

# TWCC Provider

The Taiwan Computing Cloud (TWCC) Provider is used to interact with the many resources supported by TWCC.

User should get API Key before using twcc provider.

The latest version is v0.1.0

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    twcc = {
      source  = "twcc/twcc"
      version = "0.1.0"
    }
  }
}

# Configure the TWCC Provider
provider "twcc" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}

# Example resource configuration
resource "twcc_resourcename" "example" {
  # ...
}
```

Terraform 0.12 and earlier:

```hcl
# Configure the TWCC Provider
provider "twcc" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}

# Example resource configuration
resource "twcc_resourcename" "example" {
  # ...
}
```

## Authentication

The TWCC provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static Credentials

Static credentials can be provided by adding an `apikey` and `apigw_url`
in-line in the TWCC provider block:

Usage:

```hcl
provider "twcc" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}
```

### Environment Variables

You can provide your credentials via the `TWCC_APIKEY` and
`APIGW_URL`, environment variables, representing your TWCC
API Key and APIGW URL.

```hcl
provider "twcc" {}
```

Usage:

```sh
$ export TWCC_APIKEY="yourapikey"
$ export APIGW_URL="https://apigateway.twcc.ai/"
$ terraform plan
```

## Argument Reference

* `apikey` - (Optional) This is the TWCC access key. It must be provided, but
  it can also be sourced from the `TWCC_APIKEY` environment variable.

* `apigw_url` - (Optional) This is the TWCC API gateway url. It must be provided, but
  it can also be sourced from the `APIGW_URL` environment variable.
