---
subcategory: "Firewall"
layout: "twcc"
page_title: "TWCC: twcc_firewall"
description: |-
  Provides details about a TWCC Firewall.
---

# Data Source: twcc_firewall

Use this data source to get information on an existing firewall.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_firewall" "firewall1" {
    name = "geminitestfirewall"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the firewall.

* `platform` - (Required) The name of the platform where firewall is.

* `project` - (Required) The ID of the project where firewall is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `desc` - The description of the firewall.

* `id` - The ID of the firewall.

* `status` - The status of the firewall.

* `user` - The user information who create the firewall.
