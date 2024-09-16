### Dependency

Remove DB & Middle ware

```bash
$ rm db-*
$ rm redis-*
```

`Chart.yaml`

```yaml
dependencies:
  - name: redis
    version: "17.16.0"
    repository: "oci://registry-1.docker.io/bitnamicharts"
    condition: redis.enabled
    tags:
      - middleware
  - name: postgresql-ha
    version: "11.9.0"
    repository: "oci://registry-1.docker.io/bitnamicharts"
    condition: postgresql-ha.enabled
    tags:
      - middleware
```

```bash
# download dependent charts to charts/
$ helm dependencies update
```

