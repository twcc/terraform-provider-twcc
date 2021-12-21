---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_vpn"
description: |-
  Provides details about a TWCC VPN.
---

# Data Source: twcc_vpn

Use this data source to get information on an existing VPN.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_vpn" "vpn1" {
    name = "geminitestvpn"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the VPN.

* `platform` - (Required) The name of the platform where VPN is.

* `ike_policy` - (Optional) The ID of the IKE policy whitch VPN use.

* `ipsec_policy` - (Optional) The ID of the IPsec policy whitch VPN use.

* `private_network` - (Optional) The ID of the private network whitch VPN attached to.

* `project` - (Optional) The ID of the project where VPN is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPN.

* `local_address` - The local address of the VPN.

* `local_cidr` - The local CIDR of the VPN.

* `status` - The status of the VPN.

* `user` - The user information who create the VPN.

* `vpn_connection` - The vpn connection information of the VPN.
