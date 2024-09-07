### [Tecent COS](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/resources/cos_bucket)

++ in `cvm.tf`

```json
# Tencent COS
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "private_sbucket" {
  bucket = "terraform-bucket-${local.app_id}"
  acl    = "private"
  encryption_algorithm = "AES256"
}
```

```bash
$ terraform apply --auto-approve
```

### Move .tfstate to COS

Comment in `cvm.tf`

```json
# Tencent COS
# data "tencentcloud_user_info" "info" {}

# locals {
#   app_id = data.tencentcloud_user_info.info.app_id
# }

# resource "tencentcloud_cos_bucket" "private_sbucket" {
#   bucket = "terraform-bucket-${local.app_id}"
#   acl    = "private"
#   encryption_algorithm = "AES256"
# }
```

++ in `cvm.tf`

```json
terraform {
  backend "cos" {
    region = "ap-guangzhou"
    bucket = "terraform-bucket-1329112499"
    prefix = "terraform/state"
    encrypt = true
  }
}
```

```bash
$ export TENCENTCLOUD_SECRET_KEY="..."
$ export TENCENTCLOUD_SECRET_ID="..."
$ terraform init
```

