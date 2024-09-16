### voting-app

- base 不包含 overlays 目录最小化 Kustomize 应用定义
- overlays/dev 提取镜像并复写 DB & Redis 连接信息，

```bash
$ tree kustomize  
kustomize
├── base
│   ├── kustomization.yaml
│   ├── result-deployment.yaml
│   ├── result-service.yaml
│   ├── vote-deployment.yaml
│   ├── vote-service.yaml
│   └── worker-deployment.yaml
└── overlays
    └── dev
        ├── db-deployment.yaml
        ├── db-service.yaml
        ├── kustomization.yaml
        ├── redis-deployment.yaml
        └── redis-service.yaml
```

`base/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# 引用资源
resources:
  - result-deployment.yaml
  - result-service.yaml
  - vote-deployment.yaml
  - vote-service.yaml
  - worker-deployment.yaml
```

`overlays/dev/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# 引用资源
resources:
  - ../../base
  - ./db-deployment.yaml
  - ./db-service.yaml
  - ./redis-deployment.yaml
  - ./redis-service.yaml
```

Rendering

```bash
$ kubectl kustomize base
$ kubectl kustomize overlays/dev
```

Deploy

```bash
$ kubectl kustomize . | kubectl apply -f -
```

Production using helm charts & patch

```bash
$ tree overlays/prod 
overlays/prod
└── kustomization.yaml
```

`overlays/prod/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

patches:
  - patch : |-
      - op: replace
        path: /spec/tempalte/spec/containers/0/image
        value: lyzhang19999/vote:env
    target:
      kind: Deployment
      name: vote

helmCharts:
  - name: redis
    version: "17.16.0"
    repo: "https://charts.bitnami.com/bitnami"
    valuesInline:
      fullnameOverride: redis
      auth:
        enabled: false
  - name: postgresql-ha
    version: "11.9.0"
    repo: "https://charts.bitnami.com/bitnami"
    valuesInline:
      fullnameOverride: db
      global:
        postgressql:
          username: postgres
          password: postgres
          database: postgres
          repmgrUsername: postgres
          repmgrPassword: postgres
```

Rendering

```bash
$ kubectl kustomize overlays/prod --enable-helm
```

