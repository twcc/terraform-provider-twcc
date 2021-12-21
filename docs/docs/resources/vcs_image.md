---
subcategory: "VCS"
layout: "twcc"
page_title: "TWCC: twcc_vcs_image"
description: |-
  Provides a VCS image.
---

# Resource: twcc_vcs_image

Provides a VCS image.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_solution" "testSolution" {
    name = "Ubuntu"
    project = data.twcc_project.testProject.id
    category = "os"
}

resource "twcc_vcs" "vcs1" {
    extra_property = {
        flavor = "02_vCPU_016GB_MEM_100GB_HDD"
        floating-ip = "nofloating"
        image = "Ubuntu 16.04"
        keypair = "edison63"
        private-network = "default_network"
        system-volume-type = "local_disk"
    }

    name = "geminitestvcs1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.testSolution.id
}

resource "twcc_vcs_image" "snapshot1" {
    name = "geminitestserversnap1"
    os = "Linux"
    os_version = "Ubuntu 16.04"
    platform = data.twcc_project.testProject.platform
    server = twcc_vcs.vcs1.servers[0].id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the image.

* `os` - (Required) The name of the image.

* `os_version` - (Required) The name of the image.

* `platform` - (Required) The name of the platform where image is.

* `server` - (Required) The ID of the server whitch image create from.

* `desc` - (Optional) The description of the image.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the image.

* `id` - The ID of the image.

* `is_enabled` - `True` if image is enabled. Must be `True`.

* `is_public` - `True` if image is uploaded by system admin. Must be `False`.

* `ref_img_id` - The OpenStack image ID of the image.

* `status` - The status of the image.

* `status_reason` - The status reason of the image.
