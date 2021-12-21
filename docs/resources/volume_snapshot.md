---
subcategory: "Volume"
layout: "twcc"
page_title: "TWCC: twcc_volume_snapshot"
description: |-
  Provides a volume snapshot.
---

# Resource: twcc_volume_snapshot

Provides a volume snapshot.

## Example Usage

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

resource "twcc_volume_snapshot" "snapshot" {
    name = "geminitestsnap1"
    platform = data.twcc_project.testProject.platform
    volume = twcc_volume.volume1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the snapshot.

* `platform` - (Required) The name of the platform where snapshot is.

* `volume` - (Required) The ID of the volume whitch snapshot create from.

* `desc` - (Optional) The description of the snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the snapshot.

* `id` - The ID of the snapshot.

* `restore_volume` - The ID of the volume whitch created from the snapshot.

* `snapshot_uuid` - The OpenStack ID of the snapshot.

* `status` - The status of the snapshot.

* `status_reason` - The status reason of the snapshot.

* `user` - The user information who create the snapshot.
