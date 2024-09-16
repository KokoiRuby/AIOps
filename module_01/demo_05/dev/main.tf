module "cvm" {
  source     = "../modules/cvm"
  secret_id  = var.secret_id
  secret_key = var.secret_key
  password   = var.password
  region     = var.region
}

module "k3s" {
  source     = "../modules/k3s"
  # output from module cvm
  public_ip  = module.cvm.public_ip
  private_ip = module.cvm.private_ip
}

# output from module k3s to local file
resource "local_sensitive_file" "kubeconfig" {
  content  = module.k3s.kube_config
  filename = "${path.module}/config.yaml"
}
