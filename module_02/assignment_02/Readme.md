为 Helm Demo 增加 Pre-install Hooks（Job 类型），并打印一段内容。

`module_02/assignment_02/job.yaml`

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pre-install-job
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
spec:
  template:
    spec:
      containers:
      - name: pre-install-container
        image: busybox
        command: ["sh", "-c", "echo 'helm pre-install hook'"]
      restartPolicy: Never
```

