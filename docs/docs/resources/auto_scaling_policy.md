---
subcategory: "Auto Scaling"
layout: "twcc"
page_title: "TWCC: twcc_auto_scaling_policy"
description: |-
  Provides an auto scaling policy.
---

# Resource: twcc_auto_scaling_policy

Provides an auto scaling policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
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
```

## Argument Reference

The following arguments are supported:

* `meter_name` - (Required) The monitor meter name of the auto scaling policy. Valid values: `cpu_util`, `memory.usage`, `disk.read.bytes.rate`, `disk.write.bytes.rate`, `network.incoming.bytes.rate`, `network.outgoing.bytes.rate`.

* `name` - (Required) The name of the auto scaling policy.

* `platform` - (Required) The name of the platform where auto scaling policy is.

* `project` - (Required) The ID of the project where auto scaling policy is.

* `scale_max_size` - (Required) The number of auto scale max size servers of the auto scaling policy.

* `scaledown_threshold` - (Required) The number of monitor scale down threshold value of the auto scaling policy.

* `scaleup_threshold` - (Required) The number of monitor scale up threshold value of the auto scaling policy.

* `description` - (Optional) The description of the auto scaling policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the auto scaling policy.

* `user` - The user information who create the auto scaling policy.
