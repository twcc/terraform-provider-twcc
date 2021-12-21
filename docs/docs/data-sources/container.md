---
subcategory: "Container"
layout: "twcc"
page_title: "TWCC: twcc_container"
description: |-
  Provides details about a TWCC Container.
---

# Data Source: twcc_container

Use this data source to get information on an existing container.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "k8s-taichung-default"
}

data "twcc_container" "container1" {
    name = "geminitestcontainer"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the container.

* `platform` - (Required) The name of the platform where container is.

* `project` - (Required) The ID of the project where container is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the container.

* `id` - The ID of the container.

* `pod` - The pod information of the container.

* `public_ip` - The public IP of the container.

* `service` - The service information of the container.

* `solution` - The ID of the solution whitch container create from.

* `status` - The status of the container.

* `user` - The user information who create the container.
