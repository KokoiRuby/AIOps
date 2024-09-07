### Module + Directory

```bash
$ tree
.
├── Readme.md
├── dev
│   ├── main.tf
│   └── variables.tf
├── modules
│   ├── cvm
│   │   ├── main.tf
│   │   ├── outputs.tf
│   │   ├── variables.tf
│   │   └── version.tf
│   └── k3s
│       ├── main.tf
│       ├── outputs.tf
│       └── variables.tf
└── test
    ├── main.tf
    └── variables.tf
```

`dev/main.tf`

```json
module "cvm" {
    source = "../modules/cvm"     # ref module
    secret_id = var.secret_id
    secret_key = var.secret_key
    password = var.password
}
module "k3s" {                    # ref module
    source = "../modules/k3s"
    public_ip = module.cvm.public_ip
    private_ip = module.cvm.private_ip
}
```

```bash
# under dev/test
$ terraform init
$ terraform apply --auto-approve
```

