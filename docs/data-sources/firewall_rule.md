---
subcategory: "Firewall"
layout: "twcc"
page_title: "TWCC: twcc_firewall_rule"
description: |-
  Provides details about a TWCC Firewall Rule.
---

# Data Source: twcc_firewall_rule

Use this data source to get information on an existing firewall rule.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_firewall_rule" "firewallRule1" {
    name = "geminitestfirewallrule"
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

* `action` - The action of the firewall rule.

* `destination_ip_address` - The destination IP address of the firewall rule.

* `destination_port` - The destination port of the firewall rule.

* `id` - The ID of the firewall rule.

* `ip_version` - The IP version of the firewall rule.

* `protocol` - The protocol of the firewall rule.

* `source_ip_address` - The source IP address of the firewall rule.

* `source_port` - The source port of the firewall rule.

* `user` - The user information who create the firewall rule.
