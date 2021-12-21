---
subcategory: "Extra Property"
layout: "twcc"
page_title: "TWCC: twcc_extra_property"
description: |-
  Provides details about a TWCC Extra Property.
---

# Data Source: twcc_extra_property

Use this data source to get extra property information to create container, vcs, or waf resources.

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

data "twcc_extra_property" "displayExtraProperties" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.testSolution.id
}
```

## Attributes Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where project is.

* `project` - (Required) The ID of the project to estimate extra property.

* `solution` - (Required) The ID of the solution to estimate extra property.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `extra_property` - The extra property for creating container, vcs, or waf resources.

* `id` - Same as solution ID.
