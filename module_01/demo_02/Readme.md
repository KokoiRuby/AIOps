### Keep secrets safe

[tecentcloud](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs) practice by ENV.

```bash
$ export TENCENTCLOUD_SECRET_ID="..."
$ export TENCENTCLOUD_SECRET_KEY="..."
$ terraform apply --auto-approve
```

### Import

将（不由 Terraform 管理的）资源纳入/导入 Terraform 的管理。

`diff{.tfstate, .tfstate.bak}` TENCENTCLOUD_SECRET_* info not included anymore.

```bash
$ mv terraform.tfstate terraform.tfstate.bak
# import resource with resource id 
$ terraform import tencentcloud_instance.web ins-eaiug6lk
```

```bash
# cvm.tf
resource "tencentcloud_instance" "web" {
  ...
}
```

### [Terraform cloud](https://app.terraform.io/app/organizations)

"Organization" → "Workspace" → "Create a workspace" → "CLI-Drive Workflow"

++ in `cvm.tf`

```bash
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
  cloud { 
    organization = "AnSoyo" 

    workspaces { 
      name = "Terraform" 
    } 
  } 
}

# add secret_id & secret_key back
provider "tencentcloud" {
  secret_id  = "..."
  secret_key = "..."
  region     = "ap-guangzhou"
}
```

```bash
# redirect to https://app.terraform.io/app/settings/tokens?source=terraform-login & create token given expiration
$ terraform login
# transfer .tfstate to remote
$ terraform init --auto-approve
$ terraform apply
```

### [AWS S3 Storage](https://aws.amazon.com/cn/s3/)