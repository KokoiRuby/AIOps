使用 YAML to Infra 模式创建云 Redis 数据库。

`module_01/assignment_02/tecent/redis/instance.yaml`

```yaml
# https://marketplace.upbound.io/providers/crossplane-contrib/provider-tencentcloud/v0.8.1/resources/redis.tencentcloud.crossplane.io/Instance/v1alpha1
apiVersion: redis.tencentcloud.crossplane.io/v1alpha1
kind: Instance
metadata:
  name: example-redis
spec:
  forProvider:
    passwordSecretRef:
      name: example-redis-cred
      namespace: default
      key: password
    availabilityZone: "ap-hongkong-2"
    memSize: 256
    typeId: 15
```

`module_01/assignment_02/tecent/redis/secret.yaml`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-redis-cred
type: Opaque
stringData:
  password: "Redis@123"
```

