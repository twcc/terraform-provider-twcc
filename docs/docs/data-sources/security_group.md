---
subcategory: "Security Group"
layout: "twcc"
page_title: "TWCC: twcc_security_group"
description: |-
  Provides details about a TWCC Security Group.
---

# Data Source: twcc_security_group

Use this data source to get information on an existing security group.

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
```

## Attributes Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where security group is.

* `vcs` - (Required) The ID of the vcs or waf resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the security group.

* `id` - The ID of the security group.

* `security_group_rules` - The security group rules of the security group.
