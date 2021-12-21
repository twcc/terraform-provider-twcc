---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_ipsec_policy"
description: |-
  Provides details about a TWCC IPsec Policy.
---

# Data Source: twcc_ipsec_policy

Use this data source to get information on an existing IPsec policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_ipsec_policy" "ipsec_policy1" {
    name = "geminitestipsecp"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the IPsec policy.

* `platform` - (Required) The name of the platform where IPsec policy is.

* `project` - (Required) The ID of the project where IPsec policy is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auth_algorithm` - The auth algorithm of the IPsec policy.

* `encapsulation_mode` - The encapsulation mode of the IPsec policy.

* `encryption_algorithm` - The encryption algorithm of the IPsec policy.

* `id` - The ID of the IPsec policy.

* `lifetime` - The lifetime of the IPsec policy.

* `pfs` - The PFS of the IPsec policy.

* `transform_protocol` - The transform protocol of the IPsec policy.

* `user` - The user information who create the IPsec policy.
