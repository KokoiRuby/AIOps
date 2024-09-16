### voting-app

- 要抽哪些变量和模板？
  - Secret & Configmap 等配置文件
- 调试？
  - helm template .
  - helm inspect values <remote_helm_chart>
- 多环境？
  - 使用不同的 values.yaml
- 生产环境下如何控制不部署中间件？
  - conditions

```bash
$ mkdir templates
$ touch Chart.yaml
$ touch value.yaml
```

`Chart.yaml`

```yaml
apiVersion: v2
name: vote
description: Kubernetes vote application

type: application
version: 0.1.0
appVersion: "0.1.0"
```

`values.yaml` extract image & tag

```yaml
vote:
  image: dockersamples/examplevotingapp_vote
  tag: latest

worker:
  image: dockersamples/examplevotingapp_worker
  tag: latest

result:
  image: dockersamples/examplevotingapp_result
  tag: latest
```

`vote-deployment.yaml`

```yaml
image: "{{ .Values.vote.image }}:{{ .Values.vote.tag }}"
```

`worker-deployment.yaml`

```yaml
image: "{{ .Values.worker.image }}:{{ .Values.worker.tag }}"
```

`result-deployment.yaml`

```yaml
image: "{{ .Values.result.image }}:{{ .Values.result.tag }}"
```

`values.yaml` middle ware control

```yaml
redis:
  deploy: false

postgres:
  deploy: false
```

`db-deployment.yaml` & `redis-deployment.yaml`

```yaml
{{- if .Values.postgres.deploy }}
...
{{- end }}
```

Deploy

```bash
$ helm install vote-helm . -n vote-helm --create-namespace
```

