---
subcategory: "Secret"
layout: "twcc"
page_title: "TWCC: twcc_secret"
description: |-
  Provides details about a TWCC Secret.
---

# Data Source: twcc_secret

Use this data source to get information on an existing secret.

## Example Usage

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "openstack-taichung-default-2"
}

data "twcc_secret" "secret2" {
    name = "bxSecret"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The name of the secret.

* `platform` - (Required) The name of the platform where secret is.

* `project` - (Required) The ID of the project where secret is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the secret.

* `desc` - The description of the secret.

* `expire_time` - The expire time (UTC) of the secret.

* `id` - The ID of the secret.

* `status` - The status of the secret.

* `user` - The user information who create the secret.
