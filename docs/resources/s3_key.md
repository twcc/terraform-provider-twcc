---
subcategory: "S3"
layout: "twcc"
page_title: "TWCC: twcc_s3_key"
description: |-
  Provides a private S3 key.
---

# Resource: twcc_s3_key

Provides a private S3 key.

## Example Usage

```hcl
data "twcc_project" "testProject" {
    name = "ENTxxxxxx"
    platform = "ceph-taichung-default"
}

resource "twcc_s3_key" "key1" {
    name = "geminitestkey"
    platform = data.apigw_project.testProject.platform
    project = data.apigw_project.testProject.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private S3 key.

* `platform` - (Required) The name of the platform where private S3 key is.

* `project` - (Required) The ID of the project where private S3 key is.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_key` - The S3 access key.

* `id` - The string with 'project_name-key_name'.

* `secret_key` - The S3 secret key.
