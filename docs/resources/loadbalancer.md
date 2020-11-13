---
subcategory: "LoadBalancer"
layout: "twcc"
page_title: "TWCC: twcc_loadbalancer"
description: |-
  Provides a loadbalancer.
---

# Resource: twcc_loadbalancer

Provides a loadbalancer.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENT108079"
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

resource "twcc_loadbalancer" "lb1" {
    listeners {
        name = "geminilbl1"
        pool_name = "geminilbp1"
        protocol = "HTTP"
        protocol_port = 80
    }

    name = "geminitestlb1"
    platform = data.twcc_project.testProject.platform
    pools {
        method = "ROUND_ROBIN"
        name = "geminilbp1"
        protocol = "HTTP"
        members {
            ip = "10.0.0.1"
        }

        members {
            ip = "10.0.0.2"
            port = 81
        }

        members {
            ip = "10.0.0.3"
            port = 82
            weight = 2
        }
    }

    private_net = twcc_network.network1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the loadbalancer.

* `platform` - (Required) The name of the platform where loadbalancer is.

* `private_net` - (Required) The ID of the network whitch loadbalancer attached to.

* `protocol` - (Required) The protocol for the resource of the loadbalancer. Valid values: `HTTP`, `HTTPS`, `TCP`.

* `protocol_port` - (Required) The protocol port number for the resource of the loadbalancer.

* `desc` - (Optional) The description of the loadbalancer.

* `members` - (Optional) The member resource list of the loadbalancer.

* `monitor` - (Optional) The monitor object of the loadbalancer.

The following arguments are updatable:

* `listeners` - Can not only update exist listener arguments but also add or delete listener.

* `pools` - Can not only update exist pool arguments but also add or delete pool.

### Listener Argument Reference

* `name` - (Required) The name of the listener.

* `pool_name` - (Required) The pool name of the pool whitch specified for listener attaching.

* `protocol` - (Required) The protocol for the resource of the listener. Valid values: `HTTP`, `HTTPS`, `TCP`, `TERMINATED_HTTPS`.

* `protocol_port` - (Required) The protocol port number for the resource of the listener.

* `default_tls_container_ref` - (Optional) The ID of the secret, for `TERMINATED_HTTPS` protocol listener. Must be specified if protocol is `TERMINATED_HTTPS`.

* `sni_container_refs` - (Optional) The ID of the secret list for `TERMINATED_HTTPS` protocol listener with SNI.

The following arguments are updatable:

* `default_tls_container_ref`

* `pool_name` - Change the pool that this listener uses.

* `sni_container_refs`

### Pool Argument Reference

* `method` - (Required) The load balancing algorithm of the pool. Valid values: `ROUND_ROBIN`, `LEAST_CONNECTIONS`, `SOURCE_IP`.

* `name` - (Required) The name of the pool.

* `protocol` - (Required) The protocol for the resource of the pool. Valid values: `HTTP`, `HTTPS`, `TCP`.

* `members` - (Optional) The member list information of the pool.

* `monitor` - (Optional) The monitor information of the pool. It is an one size list with map element.

The following arguments are updatable:

* `members` -  Can not only update exist member arguments but also add or delete member.

* `method`

#### Member Argument Reference

* `ip` - (Required) The IP address of the member.

* `port` - (Optional) The listening port of the member. Default is `80`.

* `weight` - (Optional) The weight value of the member. Default is `1`.

The following arguments are updatable:

* `port`

* `weight`

#### Monitor Argument Reference

* `delay` - (Optional) The delay time of the monitor. Default is `5`.

* `expected_codes` - (Optional) The expected codes (string) of the monitor.

* `http_method` - (Optional) The http method of the monitor.

* `max_retries` - (Optional) The max retries of the monitor. Default is `3`.

* `monitor_type` - (Optional) The monitor_type of the monitor. Default is `PING`.

* `timeout` - (Optional) The timeout time of the monitor. Default is `5`.

* `url_path` - (Optional) The url_path of the monitor.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `active_connections` - The active connections of the loadbalancer.

* `create_time` - The create time (UTC) of the loadbalancer.

* `id` - The ID of the loadbalancer.

* `status` - The status of the loadbalancer.

* `status_reason` - The status reason of the loadbalancer.

* `total_connections` - The total connections of the loadbalancer.

* `user` - The user information who create the loadbalancer.

* `vip` - The public IP of the loadbalancer.

* `waf` - The waf information of the loadbalancer.
