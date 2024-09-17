module "cvm" {
  source     = "../modules/cvm"
  secret_id  = var.secret_id
  secret_key = var.secret_key
  user       = var.user
  password   = var.password
  region     = var.region
}