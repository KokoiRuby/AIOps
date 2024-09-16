**Microservice**: Application src code → Dockerfile → K8s YAML manifests :cry: 数量上去不易维护管理

:smile: 使用 [Helm](https://helm.sh/) 定义应用

- 将多个微服务的工作负载、配置对象等封装为一个应用
- 屏蔽终端用户使用的复杂度
- **参数化、模板化**
- 支持多环境

**Core Concepts**

- **Chart** K8s 应用安装包，包含应用的 K8s 对象用于创建实例
- **Release** 使用默认或特定参数安装的 Helm 实例
- **Repository** 用于存储和分发的 Helm Chart 仓库 (Git/OCI)

**[CMD](https://helm.sh/docs/helm/)**

- `install` 安装 Helm Chart
- `uninstall` 卸载 Helm Release
- `get/status/list` 获取 Helm Release 的信息
- `repo add/list/remove/index` 添加/列出/移除/生成 Helm 仓库索引文件（含了仓库中可用的 Chart 列表和元数据信息）
- `search` 在 Repository 中查找 Helm Chart
- `create/package` 创建和打包 Helm Chart
- `pull` 拉取 helm chart
- `template [--debug]` 本地渲染模板文件
- `dependencies update` 更新依赖并下载 helm chart

**Advanced**

- [built-in](https://helm.sh/docs/chart_template_guide/builtin_objects/)
  - Release.Namespace
  - Values
  - Files.Get
  - Capabilities.APIVersions
- [functions](https://helm.sh/docs/chart_template_guide/function_list/)
  - String
- [dependency](https://helm.sh/docs/helm/helm_dependency/) declared in `Chart.yaml`
  - https://
  - file://
- [hooks](https://helm.sh/docs/topics/charts_hooks/)
  - `pre-install` 渲染模板之后、创建任何资源之前执行 (helm install)
  - `post-install` 所有资源部署到 Kubernetes 后执行 (helm install)
  - `pre-delete` 在从 Kubernetes 中删除任何资源之前执行 (helm uninstall)
  - `post-delete` 在删除所有资源后执行 (helm uninstall)
  - `pre-upgrade` 在渲染模板后、更新资源前执行 (helm upgrade)
  - `post-upgrade` 在升级所有资源后执行 (helm upgrade)
  - `pre-rollback` 在渲染模板之后、回滚任何资源之前执行 (helm rollback)
  - `post-rollback` 在修改所有资源后执行 (helm rollback)
  - `test` 在执行 helm test 命令时执行
  - 需要在 manifest annotations 中声明
    - `"helm.sh/hook"`
    - `"helm.sh/hook-weight"` 低值高优先
    - `"helm.sh/hook-delete-policy"` 删除策略

**3rd-Party**

- [bitnami](https://github.com/bitnami/charts)
- [artifacthub](https://artifacthub.io/)

**Sub Chart**

- dependencies 是借助 sub chart 来实现的
- 子 chart 约定放在 charts 目录，和 templates 目录同级
- `helm dependencies update` 会将依赖下载到 charts 目录
- 子 chart 也可以自己创建
- 父 chart 可以覆写子 chart 的 values.yaml 值，子不能覆盖父
- 通过 `helm inspect values <chart>` 查看哪些可以覆盖
- `global` 保留关键字，父子都可访问

**Cheatsheet**

```bash
# 安装 Helm Chart 或者执行升级，无需担心重复执行问题
$ helm upgrade --install <release name> --values <values file> <chart directory>

# 查看 Helm Chart values.yaml 配置信息
$ helm inspect values <CHART>

# 查看 Helm Repository Chart 列表
$ helm repo list

# 查看所有命名空间的 release
$ helm list --all-namespaces

# 升级应用前先更新依赖
$ helm upgrade <release> <chart> --dependency-update

# 回滚应用
$ helm history <release>
$ helm rollback <release> <revision>
```

