---
subcategory: "Auto Scaling"
layout: "twcc"
page_title: "TWCC: twcc_auto_scaling_relation"
description: |-
  Provides an auto scaling relation.
---

# Resource: twcc_auto_scaling_relation

Provides an auto scaling relation.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_vcs" "vcs1" {
    name = "geminitestvcs"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_auto_scaling_policy" "asp1" {
    meter_name = "cpu_util"
    name = "geminitestasp"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    scale_max_size = 2
    scaledown_threshold = 10
    scaleup_threshold = 50
}

resource "twcc_auto_scaling_relation" "asr1" {
    auto_scaling_policy = twcc_auto_scaling_policy.asp1.id
    platform = data.twcc_project.testProject.platform
    server = twcc_vcs.vcs1.servers[0].id
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_policy` - (Required) The ID of the auto scaling policy.

* `platform` - (Required) The name of the platform where auto scaling policy is.

* `server` - (Required) The ID of the server whitch auto scaling policy attached to.

* `loadbalancer` - (Optional) The ID of the loadbalancer whitch auto scaling policy and server attached to.

* `protocol_port` - (Optional) The number of loadbalancer protocol port (0 ~ 65535) if loadbalancer is defined.

* `scaledown_action` - (Optional) a URL whitch notification send to when auto scaling server scaledown occurs. Only support `http` protocol.

* `scaleup_action` - (Optional) a URL whitch notification send to when auto scaling server scaleup occurs. Only support `http` protocol.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The string with 'server_id/auto_scaling_policy_id'.

* `status` - The status of the auto scaling relation.

* `status_reason` - The status reason of the auto scaling relation.
