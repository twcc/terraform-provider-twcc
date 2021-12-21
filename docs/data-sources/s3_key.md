---
subcategory: "S3"
layout: "twcc"
page_title: "TWCC: twcc_s3_key"
description: |-
  Provides details about a TWCC S3 Key.
---

# Data Source: twcc_project

Use this data source to get information on an existing s3 key.

## Example Usage

### Find a Private S3 key with its Name

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "ceph-taichung-default"
}

data "twcc_s3_key" "key1" {
    name = "geminitestkey1"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

### Find a Public S3 key of the Project

```hcl
data "twcc_project" "testProject" {
  name = "ENTxxxxxx"
  platform = "ceph-taichung-default"
}

data "twcc_s3_key" "key2" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
}
```

## Attributes Reference

The following arguments are supported:

* `platform` - (Required) The name of the platform where S3 key is.

* `project` - (Required) The name of the project where S3 key is.

* `name` - (Optional) The name of the private S3 key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_key` - The access key of the S3 key.

* `id` - The string with 'project_name-key_name'.

* `is_public` - True if the key is a public S3 key.

* `secret_key` - The secret key of the S3 key.
