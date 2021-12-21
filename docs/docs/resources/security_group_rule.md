---
subcategory: "Security Group"
layout: "twcc"
page_title: "TWCC: twcc_security_group_rule"
description: |-
  Provides a security group rule.
---

# Resource: twcc_security_group_rule

Provides a security group rule.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_vcs" "site1" {
    name = "geminitestsite1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

data "twcc_security_group" "site1_sg" {
    platform = data.twcc_project.testProject.platform
    vcs = data.twcc_vcs.site1.id
}

resource "twcc_security_group_rule" "site1_sg_rule1" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    security_group = data.twcc_security_group.site1_sg.id
    direction = "egress"
    protocol = "udp"
    remote_ip_prefix = "192.168.0.0/16"
    port_range_min = 8000
    port_range_max = 8010
}
```

## Argument Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where security group rule is.

* `project` - (Required) The ID of the project where security group rule is.

* `direction` - (Optional) The direction of the security group rule. Valid values: `ingress`, `egress`. Default is `ingress`.

* `protocol` - (Optional) The protocol of the security group rule. Valid values: `ipv6-route`, `udp`, `ipv6-nonxt`, `esp`, `pgm`, `ah`, `udplite`, `egp`, `vrrp`, `tcp`, `ipv6-frag`, `ipv6-icmp`, `gre`, `ipv6-encap`, `ipv6-opts`, `rsvp`, `sctp`, `icmp`, `ospf`, `dccp`, `igmp`, `icmpv6`. Default is `tcp`.

* `port_range_max` - (Optional) The number of port range end value of the security group rule. Default is null.

* `port_range_min` - (Optional) The number of port range begin value of the security group rule. Default is null.

* `remote_ip_prefix` - (Optional) The remote IP CIDR of the security group rule. Default is `0.0.0.0/0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ethertype` - The ethertype of the security group rule.

* `id` - The ID of the security group rule.
