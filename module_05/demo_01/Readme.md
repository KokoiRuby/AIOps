Create Dockerfile

```bash
$ goctl docker -go main.go
```

Build image

```bash
$ docker build -t module_05_demo_01:latest .
```

Load image to kind cluster

```bash
$ kind load docker-image module_05_demo_01:latest --name <cluster name>
```

Apply

```bash
$ kubectl apply -f deployment.yaml rbac.yaml
```

Check

```bash
$ kubectl logs <pod>
```

