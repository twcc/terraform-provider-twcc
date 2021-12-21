---
subcategory: "VCS"
layout: "twcc"
page_title: "TWCC: twcc_vcs"
description: |-
  Provides details about a TWCC VCS.
---

# Data Source: twcc_vcs

Use this data source to get information on an existing VCS.

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
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the VCS.

* `platform` - (Required) The name of the platform where VCS is.

* `project` - (Required) The ID of the project where VCS is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the VCS.

* `id` - The ID of the VCS.

* `public_ip` - The public IP of the VCS.

* `servers` - The server list information of the VCS.

* `solution` - The ID of the solution whitch VCS create from.

* `status` - The status of the VCS.

* `user` - The user information who create the VCS.
