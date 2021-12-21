# Terraform APIGW provider plugin

## 編譯

1. test build

```unix
go build -o terraform-provider-twcc
```

2. release build

```unix
env GOOS=<Target Operating System> GOARCH=<Target Architecture> go build -o <Release Path>/terraform-provider-twcc
```

3. enable provider if terraform was installed

* windows folder

```
%APPDATA%\terraform.d\plugins
```

* linux folder

```
~/.terraform.d/plugins
```

## Terraform CLI

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    twcc = {
      source  = "twcc/twcc"
      version = "~> 1.0.0"
    }
  }
}

# Configure the TWCC Provider
provider "twcc" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}

# Example resource configuration
resource "twcc_resourcename" "example" {
  # ...
}
```

Terraform 0.12 and earlier:

```hcl
# Configure the TWCC Provider
provider "twcc" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}

# Example resource configuration
resource "twcc_resourcename" "example" {
  # ...
}
```

### debug mode

```
export TF_LOG=DEBUG
```
