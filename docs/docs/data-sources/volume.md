---
subcategory: "Volume"
layout: "twcc"
page_title: "TWCC: twcc_volume"
description: |-
  Provides details about a TWCC Volume.
---

# Data Source: twcc_volume

Use this data source to get information on an existing volume.

## Example Usage

### Find the Volume with Project ID and its Name

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_volume" "volume1" {
    name = "geminitestvolume"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

### Find the Volume with VCS or WAF Resource ID whitch Volume Attached to

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

data "twcc_volume" "volume1" {
    platform = data.twcc_project.testProject.platform
    vcs = data.twcc_vcs.vcs1.id
}
```

## Attributes Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where volume is.

* `name` - (Optional) The name of the volume.

* `project` - (Optional) The ID of the project where volume is.

* `vcs` - (Optional) The ID of the VCS or WAF whitch volume attached to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attached_host` - The attached host information of the volume.

* `create_time` - The create time (UTC) of the volume.

* `id` - The ID of the volume.

* `is_attached` - True if the volume is attached to any server.

* `is_bootable` - True if the volume is a system volume.

* `size` - The size of the volume.

* `status` - The status of the volume.

* `user` - The user information who create the volume.

* `volume_type` - The volume type of the volume.
