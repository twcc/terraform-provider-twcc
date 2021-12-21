---
subcategory: "Network"
layout: "twcc"
page_title: "TWCC: twcc_network"
description: |-
  Provides a network.
---

# Resource: twcc_network

Provides a network.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_network" "network1" {
    cidr = "10.0.0.0/24"
    gateway = "10.0.0.254"
    name = "geminitestnet1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    with_router = true
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The IP CIDR of the network.

* `gateway` - (Required) The gateway address of the network.

* `name` - (Required) The name of the network.

* `platform` - (Required) The name of the platform where network is.

* `project` - (Required) The ID of the project where network is.

* `dns_domain` - (Optional) The DNS domain name of the network.

* `nameservers` - (Optional) The nameserver list of the network.

* `with_router` - (Optional) `True` if create the network with router to connect to external network. Default is `False`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the network.

* `ext_net` - The external network name whitch network connect to if create the network with router.

* `firewall` - The firewall information whitch attached to the network.

* `id` - The ID of the network.

* `ip_version` - The IP version of the network.

* `status` - The status of the network.

* `status_reason` - The status reason of the network.

* `user` - The user information who create the network.
