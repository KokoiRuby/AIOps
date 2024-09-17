为这段 Golang [代码](https://gist.github.com/abhishekkr/3beebbc1db54b3b54914#file-tcp_server-go)写一个多阶段构建的 Dockerfile。

[goctl](https://go-zero.dev/en/docs/tutorials/cli/docker)

```bash
$ go mod init tcp_server
$ go mod tidy
$ goctl docker -go tcp_server.go
```