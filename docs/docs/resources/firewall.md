---
subcategory: "Firewall"
layout: "twcc"
page_title: "TWCC: twcc_firewall"
description: |-
  Provides a firewall.
---

# Resource: twcc_firewall

Provides a firewall.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_firewall_rule" "firewall_rule1" {
    name = "geminitestrule1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_firewall" "firewall1" {
    name = "geminitestwall1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    rules  = [twcc_firewall_rule.firewall_rule1.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the firewall.

* `platform` - (Required) The name of the platform where firewall is.

* `project` - (Required) The ID of the project where firewall is.

* `associate_networks` - (Optional) The network ID list whitch firewall attached to.

* `desc` - (Optional) The description of the firewall.

* `rules` - (Optional) The firewall rule ID list of the firewall.

The following arguments are updatable:

* `associate_networks`

* `desc`

* `rules`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the firewall.

* `id` - The ID of the firewall.

* `status` - The status of the firewall.

* `status_reason` - The status reason of the firewall.

* `user` - The user information who create the firewall.
