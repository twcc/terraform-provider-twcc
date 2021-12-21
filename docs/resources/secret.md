---
subcategory: "Secret"
layout: "twcc"
page_title: "TWCC: twcc_secret"
description: |-
  Provides a Secret.
---

# Resource: twcc_secret

Provides a Secret.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "openstack-taichung-default-2"
}

resource "twcc_secret" "secret1" {
    name = "geminitttestsecret1"
    payload = filebase64("/PATH/P12FILE.p12")
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the secret.

* `payload` - (Required) The base64 encoding certificate file payload string of the secret.

* `platform` - (Required) The name of the platform where secret is.

* `project` - (Required) The ID of the project where secret is.

* `desc` - (Optional) The description of the secret.

* `expire_time` - (Optional) The expire time (UTC) of the secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The create time (UTC) of the secret.

* `id` - The ID of the secret.

* `status` - The status of the secret.

* `user` - The user information who create the secret.
