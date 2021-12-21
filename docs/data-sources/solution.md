---
subcategory: "Solution"
layout: "twcc"
page_title: "TWCC: twcc_solution"
description: |-
  Provides details about a TWCC solution.
---

# Data Source: twcc_solution

Use this data source to get information on an existing solution.

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
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the solution.

* `project` - (Required) The ID of the project to check solution is accessible.

* `category` - (Required) The category of the solution.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `desc` - The description of the solution.

* `create_time` - The create time (UTC) of the solution.

* `id` - The ID of the solution.

* `is_public` - True if the solution is public.

* `is_tenant_admin_only` - True if the solution is only for project leader using.
