### Create a CVM

```bash
# cred in output to access cvm
$ terraform apply --auto-approve
```

### Namespace

```bash
$ sudo docker run -d nginx:latest

# get container pid by container id
$ sudo docker ps
$ sudo docker inspect --format {{.State.Pid}} container_id
$ ps -aux | grep <pid>

# get ns of proc
$ sudo lsns --task <pid>

# into mnt ns
$ sudo nsenter --target 27725 --mount

$ ls && whoami

# modify container hostname, it will modify hostname of host as well
# as we did not enter utc ns (to isolate hostname)
$ hostname my-container1 && exit && hostname

# into all ns
$ sudo nsenter --target 27725

# modify again
$ hostname my-container2 && exit && hostname
```

