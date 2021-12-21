---
subcategory: "Volume"
layout: "twcc"
page_title: "TWCC: twcc_volume"
description: |-
  Provides a volume.
---

# Resource: twcc_volume

Provides a volume.

## Example Usage

### Create Volume Normal Case

Argument `project` and `size` are also required.

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_volume" "volume1" {
    name = "geminitestvol1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    size = 1
}
```

### Create Volume by Volume Snapshot Case

Argument `src_snapshot` is also required.

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_volume_snapshot" "snapshot" {
    name = "geminitestsnapshot"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}

resource "twcc_volume" "volume2" {
    name = "geminitestvol2"
    platform = data.twcc_project.testProject.platform
    src_snapshot = data.twcc_volume_snapshot.snapshot.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.

* `platform` - (Required) The name of the platform where volume is.

* `desc` - (Optional) The description of the volume.

* `project` - (Optional) The ID of the project where volume is.

* `size` - (Optional) The number of size value of the volume.

* `src_snapshot` - (Optional) The ID of the source snapshot whitch volume created from.

* `volume_type` - (Optional) The volume type of the volume.

The following arguments are updatable:

* `size`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attached_host` - The attached host information of the volume.

* `create_time` - The create time (UTC) of the volume.

* `id` - The ID of the volume.

* `is_attached` - `True` if the volume is attaching to a server.

* `is_bootable` - `True` if the volume is a system volume.

* `is_public` - `True` if the volume is public. Must be `True`.

* `mountpoint` - The mountpoint information of the volume.

* `snapshot_list` - The snapshot list information whitch created from the volume.

* `status` - The status of the volume.

* `status_reason` - The status reason of the volume.

* `user` - The user information who create the volume.

* `volume_uuid` - The OpenStack ID of the volume.
