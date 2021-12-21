---
subcategory: "Firewall"
layout: "twcc"
page_title: "TWCC: twcc_firewall_rule"
description: |-
  Provides a firewall rule.
---

# Resource: twcc_firewall_rule

Provides a firewall rule.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the firewall rule.

* `platform` - (Required) The name of the platform where firewall rule is.

* `project` - (Required) The ID of the project where firewall rule is.

* `action` - (Optional) The action of the firewall rule. Valid values: `allow`, `deny`, `reject`. Default is `deny`.

* `destination_ip_address` - (Optional) The destination ip address or CIDR of the firewall rule.

* `destination_port` - (Optional) The destination port or port range (80:90) of the firewall rule.

* `protocol` - (Optional) The protocol of the firewall rule. Valid values: `icmp`, `tcp`, `udp`. Default is `tcp`.

* `source_ip_address` - (Optional) The source ip address or CIDR of the firewall rule.

* `source_port` - (Optional) The source port or port range (80:90) of the firewall rule.

The following arguments are updatable:

* `action`

* `destination_ip_address`

* `destination_port`

* `protocol`

* `source_ip_address`

* `source_port`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the firewall rule.

* `id` - The ID of the firewall rule.

* `ip_version` - The IP version of the firewall rule.
