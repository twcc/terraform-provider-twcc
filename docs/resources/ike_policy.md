---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_ike_policy"
description: |-
  Provides a IKE policy.
---

# Resource: twcc_ike_policy

Provides a IKE policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_ike_policy" "ike1" {
    name = "geminitestike1"
    platform = data.apigw_project.testProject.platform
    project = data.apigw_project.testProject.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the IKE policy.

* `platform` - (Required) The name of the platform where IKE policy is.

* `project` - (Required) The ID of the project where IKE policy is.

* `auth_algorithm` - (Optional) The auth algorithm of the IKE policy. Valid values: `sha1`, `sha256`, `sha384`, `sha512`. Default is `sha1`.

* `encryption_algorithm` - (Optional) The encryption algorithm of the IKE policy. Valid values: `aes-128`, `aes-192`, `aes-256`, `3des`. Default is `aes-128`.

* `ike_version` - (Optional) The IKE version of the IKE policy. Valid values: `v1`, `v2`. Default is `v1`.

* `lifetime` - (Optional) The lifetime of the IKE policy. Default is `3600`.

* `pfs` - (Optional) The perfect forward secrecy of the IKE policy. Valid values: `group2`, `group5`, `group14`. Default is `group5`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IKE policy.

* `user` - The user information who create the IKE policy.
