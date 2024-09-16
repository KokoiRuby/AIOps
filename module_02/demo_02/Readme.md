### Examine

```bash
$ docker build -t nginx:test -f Dockerfile1 .
$ mkdir -p /tmp/module_02/demo_02
$ docker save nginx:test -o /tmp/module_02/demo_02/nginx.tar
$ cd /tmp/module_02/demo_02
$ extract nginx.tar
```

### Overlay

```bash
# into cvm
$ cd /tmp/overlay

# lower_dir is read
# upper_dir is write-able
# merged_dir is view
# work_dir is middle layer
$ ls lower_dir && ls upper_dir && ls merged_dir && ls work_dir

# create overlay
$ sudo mount -t overlay \
	-o lowerdir=lower_dir/,upperdir=upper_dir/,workdir=work_dir/ \
	none merged_dir/

# chk
$ df -a
$ ls merged_dir/

# upper overwrites lower
$ cat merged_dir/same.txt

$ vi merged_dir/1.txt
# lower not modified since read-only
$ cat lower_dir/1.txt
# ++ = write-on-copy
$ cat upper_dir/1.txt


$ rm merged_dir/2.txt
# lower not deleted since read-only
$ ls lower_dir/
# marked C = deleted
$ ls -al upper_dir/
$ ls -al merged_dir/
```

### Cache

```bash
$ docker build -t pyapp:test -f Dockerfile2 .
$ docker build -t pyapp:test -f Dockerfile3 .
```

