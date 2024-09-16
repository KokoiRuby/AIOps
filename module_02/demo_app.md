### K8s Auto-recovery

依赖服务未就绪，业务进程以非 0 状态码退出。

K8s 检测到容器异常退出，指数退避 backoff 自动重启直到依赖服务就绪 :cry: 这会导致应用整体启动的时间变长...

:smile: initContainers 等待依赖服务就绪后再启动主容器 [k8s-wait-for](https://github.com/groundnuty/k8s-wait-for) 实现链式调用控制启动顺序。

### DB Init

业务代码初始化 or K8s Job (推荐) git clone <sql_schema_repo> 然后执行 SQL 初始化数据库。

### Middleware HA

最佳实践：使用社区提供的 Helm Chart

- [Postgres](https://github.com/bitnami/charts/tree/main/bitnami/postgresql-ha)
- [Redis](https://github.com/bitnami/charts/blob/main/bitnami/redis/README.md)

生产：尽量使用云服务，自托管次之 self-hosting