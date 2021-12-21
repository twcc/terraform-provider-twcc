---
subcategory: "Container"
layout: "twcc"
page_title: "TWCC: twcc_container"
description: |-
  Provides a container.
---

# Resource: twcc_container

Provides a container.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "k8s-taichung-default"
}

data "twcc_solution" "testSolution" {
    name = "Tensorflow"
    project = data.twcc_project.testProject.id
    category = "container"
}

resource "twcc_container" "container1" {
    extra_property = {
        flavor = "1 GPU + 04 cores + 090GB memory"
        gpfs01-mount-path = "/mnt"
        gpfs02-mount-path = "/tmp"
        image = "tensorflow-19.08-py3:latest"
        replica = "1"
    }

    name = "geminitestcontainer"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.testSolution.id
}
```

## Argument Reference

The following arguments are supported:

* `extra_property` - (Required) The extra property dictionary to create a container resource.

* `name` - (Required) The name of the container.

* `platform` - (Required) The name of the platform where container is.

* `project` - (Required) The ID of the project where container is.

* `solution` - (Required) The ID of the solution whitch container create from.

* `desc` - (Optional) The description of the container.

### Get extra_property please reference `data-source/extra_property.md`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the container.

* `id` - The ID of the container.

* `pod` - The pod information of the container.

* `public_ip` - The public IP of the container.

* `service` - The service information of the container.

* `status` - The status of the container.

* `status_reason` - The status reason of the container.

* `user` - The user information who create the container.
