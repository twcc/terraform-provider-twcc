---
subcategory: "VCS"
layout: "twcc"
page_title: "TWCC: twcc_vcs"
description: |-
  Provides a VCS.
---

# Resource: twcc_vcs

Provides a VCS.

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
```

## Argument Reference

The following arguments are supported:

* `extra_property` - (Required) The extra property dictionary to create a VCS resource.

* `name` - (Required) The name of the VCS.

* `platform` - (Required) The name of the platform where VCS is.

* `project` - (Required) The ID of the project where VCS is.

* `solution` - (Required) The ID of the solution whitch VCS create from.

* `desc` - (Optional) The description of the VCS.

### Get extra_property please reference `data-source/extra_property.md`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the VCS.

* `id` - The ID of the VCS.

* `public_ip` - The public IP of the VCS.

* `servers` - The server list information of the VCS.

* `status` - The status of the VCS.

* `status_reason` - The status reason of the VCS.

* `user` - The user information who create the VCS.
