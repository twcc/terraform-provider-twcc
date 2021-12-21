---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_ipsec_policy"
description: |-
  Provides a IPsec policy.
---

# Resource: twcc_ipsec_policy

Provides a IPsec policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_ipsec_policy" "ipsec1" {
    name = "geminitestipsec1"
    platform = data.apigw_project.testProject.platform
    project = data.apigw_project.testProject.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the IPsec policy.

* `platform` - (Required) The name of the platform where IPsec policy is.

* `project` - (Required) The ID of the project where IPsec policy is.

* `auth_algorithm` - (Optional) The auth algorithm of the IPsec policy. Valid values: `sha1`, `sha256`, `sha384`, `sha512`. Default is `sha1`.

* `encapsulation_mode` - (Optional) The encapsulation mode of the IPsec policy. Valid values: `tunnel`, `transport`. Default is `tunnel`.

* `encryption_algorithm` - (Optional) The encryption algorithm of the IPsec policy. Valid values: `aes-128`, `aes-192`, `aes-256`, `3des`. Default is `aes-128`.

* `lifetime` - (Optional) The lifetime of the IPsec policy. Default is `3600`.

* `pfs` - (Optional) The perfect forward secrecy of the IPsec policy. Valid values: `group2`, `group5`, `group14`. Default is `group5`.

* `transform_protocol` - (Optional) The transform protocol of the IPsec policy. Valid values: `esp`, `ah`, `ah-esp`. Default is `esp`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IPsec policy.

* `user` - The user information who create the IPsec policy.
