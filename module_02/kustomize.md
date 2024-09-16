[Kustomize](https://kustomize.io/)

A CLI tool (already merged into kubectl) to **overwrite** manifests for multi-envs.

Base + Overlays + `kustomization.yaml` → manifest.yaml

base 以及每一个 overylay 目录中都需要包含一个 `kustomization.yaml`

```bash
├── base # 基础目录，或者不同环境的通用 Manifest
│ ├── deployment.yaml
│ ├── kustomization.yaml
└── overlays
    ├── dev                   # dev 环境目录
    │ ├── hpa.yaml
    │ └── kustomization.yaml
    ├── production            # production 环境目录
    │ ├── hpa.yaml
    │ ├── kustomization.yaml
    └── staging               # staging 环境目录
    ├── hpa.yaml
    ├── kustomization.yaml
```

[kustomization.yaml](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/)

- `resource` 定义 Manifest 资源，文件、目录或 URL
2. `secretGenerator` 生成 secret 对象
3. `configMapGenerator` 生成 configMap 对象
4. `images` 覆写 image tag
5. `helmCharts` 定义依赖的 helm chart
6. `patchesStrategicMerge` 覆写操作，v1 版本将废弃，建议迁移到 patches

**Why Kustomize?**

- 支持 override manifest 中任何字段，如增加 Annotations/Labels
- 编写简单，无需学习复杂的模板变量语法
- 多环境下、Base 和 Overlays 模式能够实现很好的 Manifest 复用
- Helm 隐藏应用细节，对最终用户友好
- Kustomize 暴露所有 K8s API，对开发者友好

:cry:

- 分发没有 Helm Chart 方便
- 生态相比较 Helm 差
- 强依赖目录结构