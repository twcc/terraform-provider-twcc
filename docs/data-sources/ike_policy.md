---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_ike_policy"
description: |-
  Provides details about a TWCC IKE Policy.
---

# Data Source: twcc_ike_policy

Use this data source to get information on an existing IKE policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_ike_policy" "ike_policy1" {
    name = "geminitestikep"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the IKE policy.

* `platform` - (Required) The name of the platform where IKE policy is.

* `project` - (Required) The ID of the project where IKE policy is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auth_algorithm` - The auth algorithm of the IKE policy.

* `encryption_algorithm` - The encryption algorithm of the IKE policy.

* `ike_version` - The IKE version of the IKE policy.

* `id` - The ID of the IKE policy.

* `lifetime` - The lifetime of the IKE policy.

* `pfs` - The PFS of the IKE policy.

* `user` - The user information who create the IKE policy.
