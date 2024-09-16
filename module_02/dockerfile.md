### Dockerfile

**History of Container**

- 1979~ Chroot 文件系统隔离
- 2000~ FreeBSD Jails 早期容器隔离
- 2002~ Linux Namespaces 进程隔离
- 2004~ Solaris Zones 快照/克隆
- 2006~ Google Process Containers
- 2007~ Linux Control Group 重命名 Process Containers 合并进 Kernel 2.6.24
- 2008~ LXC Linux Containers 使用 namespace 和 cgroupss 实现容器管理
- **2013~ Docker 使用 libcontainer 实现容器**

**What is Container?**

A **process** on host constained by **namespaces** (limits what u can see) & **cgroup** (limits how much u can use).

**vs Image?**

A Conatiner is running instance of Image. A image encapsulates everything a container to run.

**Dockerfile**

A group of CMDs to describe how to build a container image (from a base image).

**Essenstial of Image**

- OverlayFS (lower/upper/merged)
- Layer stack
- Modify = push new layer to stack
- Layer can be cached 
- Named by hash

**[CMD](https://docs.docker.com/reference/dockerfile/)**

**Cache**

Any delta on lower layer will de-cache upper layer.

Move layers that less-like to delta to as lower as possible.

### **[Best Practice](https://docs.docker.com/build/building/best-practices/)**

- As less layers as possible
- Appropriate order of CMD to re-use cache
- .dockerignore
- Base image
- Multi-stage
- Non-root
- Misc
  - Avoid latest tag
  - Official image
  - Timezone
- Ext.
  - DinD = Docker in Docker :cry: Security
  - [buildkit](https://github.com/moby/buildkit) 不依赖于 Dockerd 构建镜像
  - [kaniko](https://github.com/GoogleContainerTools/kaniko) 于 K8s 中构建镜像

```bash
# DinD
$ docker run -v /var/run/docker.sock:/var/run/docker.sock -it docker
```

:confused: **Build once, Run everywhere?**

Docker 在拉取镜像的时候会根据当前主机的架构从 Registry 拉取对应的镜像。

架构信息被包含在镜像的 manifest 中。使用 [crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane) 可以查看。

```bash
$ crane manifest alpine:latest | jq .
$ crane manifest alpine:latest@<digest> | jq .
```

```bash
$ mkdir tmp && cd tmp && crane export -v alpine:latest - | tar xv
```

**Multi platform build by [QEMU](https://docs.docker.com/build/building/multi-platform/#qemu)**

```bash
# build multi-platform
$ docker build --platform linux/amd64,linux/arm64 -t test:v1 -f Dockerfile .
```

### Container Runtime

Docker → dockerd → containerd → OCI → **runc** → container

K8s → CRI → containerd/CRI-O  → OCI → **runc** → container

runc 不能直接运行容器镜像，需要镜像中的 root filesystem + [config.json](https://github.com/opencontainers/runtime-spec/blob/main/config.md) = OCI Runtime Bundle

```bash
# in CVM
$ cd /tmp

# out-of-box rootfs
$ ls rootfs

# gen config.json
$ runc spec
$ cat config.json

# run
$ sudo runc run my-container

# chk
$ ps aux
$ hostname

# exit
$ exit
```

