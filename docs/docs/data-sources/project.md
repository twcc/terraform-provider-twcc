---
subcategory: "Project"
layout: "twcc"
page_title: "TWCC: twcc_project"
description: |-
  Provides details about a TWCC Project.
---

# Data Source: twcc_project

Use this data source to get information on an existing project.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the project.

* `platform` - (Required) The name of the platform where project is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the project.
