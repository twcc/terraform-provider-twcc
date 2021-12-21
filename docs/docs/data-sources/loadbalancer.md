---
subcategory: "Loadbalancer"
layout: "twcc"
page_title: "TWCC: twcc_loadbalancer"
description: |-
  Provides details about a TWCC Loadbalancer.
---

# Data Source: twcc_loadbalancer

Use this data source to get information on an existing loadbalancer.

## Example Usage

### Find Loadbalancer by Name and Project ID

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_loadbalancer" "lb1" {
    name = "geminitestlb1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

### Find Loadbalancer by Name and Private Network ID

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

data "twcc_loadbalancer" "lb2" {
    name = "geminitestlb2"
    platform = data.twcc_network.network1.platform
    private_net = data.twcc_network.network1.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the loadbalancer.

* `platform` - (Required) The name of the platform where loadbalancer is.

* `private_net` - (Optional) The ID of the network whitch loadbalancer attach to.

* `project` - (Optional) The ID of the project where loadbalancer is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the loadbalancer.

* `desc` - The description of the loadbalancer.

* `listeners` - The load balancer listener list information of the loadbalancer.

* `id` - The ID of the loadbalancer.

* `pools` - The load balancer pool list information of the loadbalancer.

* `status` - The status of the loadbalancer.

* `user` - The user information who create the loadbalancer.
