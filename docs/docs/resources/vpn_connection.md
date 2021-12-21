---
subcategory: "VPN"
layout: "twcc"
page_title: "TWCC: twcc_vpn_connection"
description: |-
  Provides a VPN connection.
---

# Resource: twcc_vpn_connection

Provides a VPN connection.

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

resource "twcc_vpn_connection" "vc1" {
    peer_address = "10.0.0.254"
    peer_cidrs = ["10.0.0.0/24"]
    platform = data.twcc_project.testProject.platform
    psk = "testgemini"
    vpn = data.twcc_vpn.vpn1.id
}
```

## Argument Reference

The following arguments are supported:

* `peer_address` - (Required) The peer IP address of the connection.

* `peer_cidrs` - (Required) The peer IP CIDR list of the connection.

* `platform` - (Required) The name of the platform where connection is.

* `psk` - (Required) The pre-shared key setting of the connection.

* `vpn` - (Required) The ID of the VPN for connection.

* `dpd_action` - (Optional) The dead peer detection action of the connection. Valid values: `clear`, `hold`, `restart`, `disabled`, `restart-by-peer`. Default is `hold`.

* `dpd_interval` - (Optional) The dead peer detection query and delay interval (seconds) of the connection. Default is `30`.

* `dpd_timeout` - (Optional) The dead peer detection timeout of the connection. Should be greater than DPD interval. Default is `120`.

* `initiator` - (Optional) Indicates whether this VPN can only respond to connections or both respond to and initiate connections. Valid values: `response- only`, `bi-directional`. Default is `bi-directional`.

* `mtu` - (Optional) The maximum transmission unit value to address fragmentation of the connection. Minimum value is `68` for IPv4, and `1280` for IPv6. Default is `1500`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The string with 'vpn_id-connection'.

* `peer_id` - The peer ID of the connection.

* `status` - The status of the connection.
