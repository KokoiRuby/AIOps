### [Install](https://developer.hashicorp.com/terraform/install) Terraform

Get [tecentcloud](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs) provider from [Registry](https://registry.terraform.io/browse/providers).

### AuthN TecentCloud

[Console](https://console.cloud.tencent.com/) → ”[访问管理](https://console.cloud.tencent.com/cam/overview)“ → ”用户列表“ → ”新建用户“ → ”用户详情“ → ”API 密钥“ copy to `terraform.tf`

```json
...
provider "tencentcloud" {
  secret_id  = "..."
  secret_key = "..."
  region     = "ap-guangzhou"
}
...
# Get availability instance types
data "tencentcloud_instance_types" "default" {
  cpu_core_count = 2
  # memory_size    = 1
}
...
```

```bash
$ terraform init
# Plan: 4 to add, 0 to change, 0 to destroy.
$ terraform plan
```

### Provision VM on TecentCloud

```bash
$ terraform apply
$ terraform apply -auto-approve
```

[Console](https://console.cloud.tencent.com/) → ”[云服务器](https://console.cloud.tencent.com/cvm/instance/index?rid=1)“ → Switch region if necessary

### Destroy

```bash
$ terraform destroy
$ terraform destroy -auto-approve
```

### Audit

[Console](https://console.cloud.tencent.com/) → ”[操作审计](https://console.cloud.tencent.com/cloudaudit)“