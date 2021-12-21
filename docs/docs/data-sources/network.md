---
subcategory: "Network"
layout: "twcc"
page_title: "TWCC: twcc_network"
description: |-
  Provides details about a TWCC Network.
---

# Data Source: twcc_network

Use this data source to get information on an existing network.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_network" "network1" {
    name = "geminitestnet"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the network.

* `platform` - (Required) The name of the platform where network is.

* `project` - (Required) The ID of the project where network is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cidr` - The fixed IP CIDR of the network.

* `create_time` - The create time (UTC) of the network.

* `dns_domain` - The DNS domain of the network.

* `ext_net` - The external network name whitch network connect to if create the network with router.

* `gateway` - The gateway of the network.

* `id` - The ID of the network.

* `ip_version` - The IP version of the network.

* `nameservers` - The nameserver list of the network.

* `status` - The status of the network.

* `user` - The user information who create the network.

* `with_router` - True if create the network with router to connect to external network.
