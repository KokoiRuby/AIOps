### Create VM

`version.tf`

```json
terraform {
  required_version = "> 0.13.0"
  required_providers {
    tencentcloud = {
      source  = "tencentcloudstack/tencentcloud"
      version = "1.81.5"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.14"
    }
  }
}

# to install argo cd
provider "helm" {
  kubernetes {
    # local kubeconfig
    config_path = local_sensitive_file.kubeconfig.filename
  }
}
```

`variables.tf`

```json
variable "secret_id" {
  default = "your_secret"
}

variable "secret_key" {
  default = "your_secret"
}

variable "region" {
  default = "ap-guangzhou"
}

# pwd for VM SSH login
variable "passw" {
  default = "password123"
}
```

split `main.tf` into:

`cvm.tf`

```json
# Configure the TencentCloud Provider
# Read from variables.tf
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
}

resource "tencentcloud_instance" "web" {
  ...
  password                   = var.passw
}
```

`output.tf`

```json
output "public_ip" {
    description = "vm public ip"
    value = tencentcloud_instance.web[0].public_ip
}

output "kube_config" {
    description = "kubeconfig"
    value = "${path.module}/config.yaml"
}

output "password" {
    description = "vm ssh password"
    value = var.password
}
```

Pass secre_* ENV to `variables.tf`

```bash
$ export TF_VAR_secret_id="..."
$ export TF_VAR_secret_key="..."
```

Apply

```bash
$ terraform apply --auto-approve
```

### Deploy K3s

`k3s.tf` 

[k3s module](https://registry.terraform.io/modules/xunleii/k3s/module/latest)

```json
module "k3s" {
  source                   = "xunleii/k3s/module"
  k3s_version              = "v1.28.11+k3s2"
  generate_ca_certificates = true
  global_flags = [
    "--tls-san ${tencentcloud_instance.web[0].public_ip}",
    "--write-kubeconfig-mode 644",
    "--disable=traefik",
    "--kube-controller-manager-arg bind-address=0.0.0.0",
    "--kube-proxy-arg metrics-bind-address=0.0.0.0",
    "--kube-scheduler-arg bind-address=0.0.0.0"
  ]
  k3s_install_env_vars = {}

  servers = {
    "k3s" = {
      ip = tencentcloud_instance.web[0].private_ip
      connection = {
        timeout  = "60s"
        type     = "ssh"
        host     = tencentcloud_instance.web[0].public_ip
        password = var.password
        user     = "ubuntu"
      }
    }
  }
}

resource "local_sensitive_file" "kubeconfig" {
  content  = module.k3s.kube_config
  filename = "${path.module}/config.yaml"
}

```

### Install ArgoCD via Helm

`helm.tf`

[helm module](https://registry.terraform.io/providers/hashicorp/helm/latest)

```json
resource "helm_release" "argo_cd" {
  depends_on       = [module.k3s]
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true
}
```

### Check Dependency

```bash
$ sudo apt install graphviz
$ terraform graph | dot -Tsvg > graph.svg
```

