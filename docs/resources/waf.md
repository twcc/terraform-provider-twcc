---
subcategory: "WAF"
layout: "twcc"
page_title: "TWCC: twcc_waf"
description: |-
  Provides a WAF.
---

# Resource: twcc_waf

Provides a WAF.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

data "twcc_solution" "testSolution" {
    name = "F5_WAF"
    project = data.twcc_project.testProject.id
    category = "waf"
}

resource "twcc_waf" "waf1" {
    extra_property = {
        availability-zone = "nova"
        flavor = "08_core_040GB_memory_160GB_disk"
        image = "F5-AWAF-Production"
        password = "password"
        private-network = "default_network"
    }

    name = "geminitestwaf1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    solution = data.twcc_solution.testSolution.id
}
```

## Argument Reference

The following arguments are supported:

* `extra_property` - (Required) The extra property dictionary to create a WAF resource.

* `name` - (Required) The name of the WAF.

* `platform` - (Required) The name of the platform where WAF is.

* `project` - (Required) The ID of the project where WAF is.

* `solution` - (Required) The ID of the solution whitch WAF create from.

* `desc` - (Optional) The description of the WAF.

### Get extra_property please reference `data-source/extra_property.md`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the WAF.

* `id` - The ID of the WAF.

* `public_ip` - The public IP of the WAF.

* `servers` - The server list information of the WAF.

* `status` - The status of the WAF.

* `status_reason` - The status reason of the WAF.

* `user` - The user information who create the WAF.
