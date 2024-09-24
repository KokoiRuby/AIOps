resource "helm_release" "crossplane" {
  depends_on       = [module.k3s, null_resource.download_k3s_yaml]
  name             = "crossplane"
  repository       = "https://charts.crossplane.io/stable"
  chart            = "crossplane"
  namespace        = "crossplane"
  create_namespace = true
}
