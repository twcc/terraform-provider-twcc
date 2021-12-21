---
subcategory: "Auto Scaling"
layout: "twcc"
page_title: "TWCC: twcc_auto_scaling_policy"
description: |-
  Provides details about a TWCC Auto Scaling Policy.
---

# Data Source: twcc_auto_scaling_policy

Use this data source to get information on an existing auto scaling policy.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_auto_scaling_policy" "asp1_data" {
    name = "geminitestasp2"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the auto scaling policy.

* `platform` - (Required) The name of the platform where auto scaling policy is.

* `project` - (Required) The ID of the project where auto scaling policy is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - The description of the auto scaling policy.

* `meter_name` - The monitor meter name of the auto scaling policy.

* `id` - The ID of the auto scaling policy.

* `scale_max_size` - The server auto scale max size of the auto scaling policy.

* `scaledown_threshold` - The monitor scale down threshold value of the auto scaling policy.

* `scaleup_threshold` - The monitor scale up threshold value of the auto scaling policy.

* `user` - The user information who create the auto scaling policy.
