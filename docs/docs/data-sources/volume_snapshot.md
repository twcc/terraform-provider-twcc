---
subcategory: "Volume"
layout: "twcc"
page_title: "TWCC: twcc_volume_snapshot"
description: |-
  Provides details about a TWCC Volume Snapshot.
---

# Data Source: twcc_volume_snapshot

Use this data source to get information on an existing volume snapshot.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_volume_snapshot" "snapshot1" {
    name = "geminitestsnapshot"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the volume snapshot.

* `platform` - (Required) The name of the platform where volume snapshot is.

* `project` - (Required) The ID of the project where volume snapshot is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the volume snapshot.

* `desc` - The description of the volume snapshot.

* `id` - The ID of the volume snapshot.

* `status` - The status of the volume snapshot.

* `user` - The user information who create the volume snapshot.

* `volume` - The source volume of the volume snapshot.
