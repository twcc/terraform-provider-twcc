---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_vpn"
description: |-
  Provides a VPN.
---

# Resource: twcc_vpn

Provides a VPN.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_ike_policy" "ike1" {
    name = "geminitestike1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_ipsec_policy" "ipsec1" {
    name = "geminitestipsec1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_vpn" "vpn1" {
    ike_policy = twcc_ike_policy.ike1.id
    ipsec_policy = twcc_ipsec_policy.ipsec1.id
    name = "geminitestvpn1"
    platform = data.twcc_project.testProject.platform
    private_network = twcc_network.network1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPN.

* `ike_policy` - (Required) The ID of the IKE policy of the VPN.

* `ipsec_policy` - (Required) The ID of the IPsec policy of the VPN.

* `platform` - (Required) The name of the platform where VPN is.

* `private_network` - (Required) The ID of the network whitch VPN connected to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPN.

* `local_address` - The local address of the VPN.

* `local_cidr` - The local CIDR of the VPN.

* `status` - The status of the VPN.

* `user` - The user information who create the VPN.
