---
subcategory: "WAF"
layout: "twcc"
page_title: "TWCC: twcc_waf"
description: |-
  Provides details about a TWCC WAF.
---

# Data Source: twcc_waf

Use this data source to get information on an existing WAF.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_waf" "waf1" {
    name = "geminitestwaf"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the WAF.

* `platform` - (Required) The name of the platform where WAF is.

* `project` - (Required) The ID of the project where WAF is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the WAF.

* `id` - The ID of the WAF.

* `public_ip` - The public IP of the WAF.

* `servers` - The server list information of the WAF.

* `solution` - The ID of the solution whitch WAF create from.

* `status` - The status of the WAF.

* `user` - The user information who create the WAF.
