### Practice

`--dry-run=client` to gen template

```bash
$ kubectl create deployment nginx \
	--image=nginx -o yaml \
	--dry-run=client
	
$ kubectl create service clusterip my-cs \
	--tcp=5678:8080 -o yaml \
	--dry-run=client
	
$ kubectl create configmap my-config \
	--from-literal=key1=config1 \
	--from-literal=key2=config2 \
	-o yaml --dry-run=client
	
$ kubectl create [workload|config|storage] --help
```

`---` to merge multiple manifests

```yaml
<manifest1>
---
<manifest1>
---
...
```

**Bool**

```bash
# will be transformed to bool
y|Y|yes|Yes|YES|n|N|no|No|NO|true|True|TRUE|false|False|FALSE|on|On|ON|off|Off|OFF
```

```yaml
# use ""
env:
- name: GERMANY
  value: DE
- name: NORWAY
  value: "NO"
```

Numeric

```yaml
env:
- name: GO_VERSION
  value: "1.10"
```

**Multi-lines**

`>` folded 折叠

`|` literal 字面量

`-` no trailing newline 删除换行，删除尾部空行

`+` keep trailing newline 保留换行，保留尾部空行

```yaml
# newlines are replaced with spaces.
# Your long string here.(\n)
key: >
Your long
string here.

# without a trailing newline.
# Your long string here.
key: >-
Your long
string here.

# newlines are preserved.
# ### Heading
# * Bullet
# * Points(\n)
key: |
### Heading
* Bullet
* Points

# without a trailing newline.
# ### Heading
# * Bullet
# * Points
key: |-
### Heading
* Bullet
* Points
```

Reference

```yaml
spec:
  selector:
    # def
    matchLabels: &pod-labels
      app: myapp
template:
  metadata:
    # ref
    labels: *pod-labels
```

### Tool

[He3](https://he3app.com/) YAML Viewer

VSCode YAML plugin

[yq](https://github.com/mikefarah/yq)