---
subcategory: "Volume"
layout: "twcc"
page_title: "TWCC: twcc_attachment"
description: |-
  Provides a volume attachment.
---

# Resource: twcc_volume_attachment

Provides a volume attachment.

## Example Usage

### Create Volume Normal Case

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

resource "twcc_volume" "volume1" {
    name = "geminitestvol1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    size = 1
}

resource "twcc_volume_attachment" "volume_attachment" {
    platform = data.twcc_project.testProject.platform
    server = data.twcc_vcs.vcs1.servers[0].id
    volume = twcc_volume.volume1.id
}
```

## Argument Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where volume is.

* `server` - (Required) The ID of the server whitch volume attached to.

* `volume` - (Required) The ID of the volume.

* `mountpoint` - (Optional) The mountpoint of the volume attachment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - * `id` - The string with 'server_id/volume_id'.
