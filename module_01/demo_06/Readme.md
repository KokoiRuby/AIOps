`demo_6/tencent/helm.tf`

```json
resource "helm_release" "crossplane" {
  depends_on       = [module.k3s]
  name             = "crossplane"
  repository       = "https://charts.crossplane.io/stable"
  chart            = "crossplane"
  namespace        = "crossplane"
  create_namespace = true
}
```

```bash
$ terraform apply --auto-approve
```

[provider-tecentcloud](https://marketplace.upbound.io/providers/crossplane-contrib/provider-tencentcloud/v0.8.3)

```yaml
# provider
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-tencentcloud
spec:
  package: xpkg.upbound.io/crossplane-contrib/provider-tencentcloud:v0.8.3
---
apiVersion: tencentcloud.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    secretRef:
      key: credentials
      name: example-creds
      namespace: crossplane
    source: Secret
---
apiVersion: v1
kind: Secret
metadata:
  name: example-creds
  namespace: crossplane
type: Opaque
stringData:
  credentials: |
    {
      "secret_id": "...",
      "secret_key": "...",
      "region": "ap-guangzhou"
    }
```

```yaml
# cvm
apiVersion: vpc.tencentcloud.crossplane.io/v1alpha1
kind: VPC
metadata:
  name: resource-vpc
  namespace: crossplane-system
spec:
  forProvider:
    cidrBlock: 10.1.0.0/16
    name: crossplane-test-vpc
---
apiVersion: vpc.tencentcloud.crossplane.io/v1alpha1
kind: Subnet
metadata:
  name: example-cvm-subnet
spec:
  forProvider:
    availabilityZone: "ap-guangzhou-2"
    cidrBlock: "10.2.2.0/24"
    name: "test-crossplane-cvm-subnet"
    # ref
    vpcIdRef:
      name: "example-cvm-vpc"
---
apiVersion: cvm.tencentcloud.crossplane.io/v1alpha1
kind: Instance
metadata:
  name: example-cvm
spec:
  forProvider:
    instanceName: "test-crossplane-cvm"
    availabilityZone: "ap-guangzhou-2"
    instanceChargeType: "SPOTPAID"
    imageId: "img-487zeit5"
    instanceType: "SA5.MEDIUM4"
    systemDiskType: "CLOUD_BSSD"
    # ref
    vpcIdRef:
      name: "example-cvm-vpc"
    subnetIdRef:
      name: "example-cvm-subnet"
```

```bash
$ export KUBECONFIG="$(pwd)/config.yaml"
$ kubectl apply -f provider.yaml
$ kubectl apply -f providerConfig.yaml
# edit secret.yaml and fill the secret
$ kubectl apply -f secret.yaml
$ kubectl apply -f cvm/
```

