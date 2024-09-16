```bash
$ helm create demo
$ tree
demo
├── Chart.yaml    # meta & desc
├── charts
├── templates     # manifests
│   ├── NOTES.txt
│   ├── _helpers.tpl
│   ├── deployment.yaml
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── service.yaml
│   ├── serviceaccount.yaml
│   └── tests
│       └── test-connection.yaml
└── values.yaml    # for rendering templates
```

```bash
$ helm template .
```

```bash
$ helm instal nginx . -n nginx --create-namespace --set image.tag='latest'
```

```bash
$ helm uninstall nginx
```

```bash
# install if not exists, upgrade if exists
$ helm upgrade --install vote-helm . -n vote-helm --create-namespace
```

```bash
$ helm uninstall vote-helm
```

```bash
# specify a remote helm chart repo
$ helm upgrade --install ingress-nginx ingress-nginx \
	--repo https://kubernetes.github.io/ingress-nginx \
	--namespace ingress-nginx --create-namespace
```