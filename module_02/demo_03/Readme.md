### Dockerfile1

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
WORKDIR /src
COPY main.go .
RUN go build -o /bin/example ./main.go

FROM debian:stable-slim
COPY --from=builder /bin/example /bin/example
CMD [ "/bin/example" ]
```

```bash
$ docker build -t app:test . -f Dockerfile1
$ docker run app:test --name app
```

### Dockerfile2

```bash
# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
WORKDIR /src
COPY main.go .
RUN go build -o /bin/example ./main.go

FROM alpine:3.20.2
COPY --from=builder /bin/example /bin/example
CMD [ "/bin/example" ]
```

```bash
$ docker build -t app:test . -f Dockerfile2
$ docker run app:test --name app
```

### Dockerfile3

```dockerfile
# syntax=docker/dockerfile:1
# GNU C
FROM golang:1.21 AS builder
WORKDIR /src
COPY ./example .
RUN go mod tidy && go build -o /bin/example ./main.go

# musl libc
FROM alpine:3.20.2
COPY --from=builder /bin/example /bin/example
CMD [ "/bin/example" ]
```

```bash
$ docker build -t app:test . -f Dockerfile3
# failed to run due to alpine does not contain GNU C
$ docker run app:test --name app
```

### Dockerfile4

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
WORKDIR /src
COPY ./example .
# disable CGO
RUN go mod tidy && CGO_ENABLED=0 && go build -o /bin/example ./main.go

FROM alpine:3.20.2
COPY --from=builder /bin/example /bin/example
CMD [ "/bin/example" ]
```

```bash
$ docker build -t app:test . -f Dockerfile4
# ok
$ docker run app:test --name app
```

### Dockerfile5

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
WORKDIR /src
COPY ./example .
# disable CGO
RUN go mod tidy && CGO_ENABLED=0 && go build -o /bin/example ./main.go

FROM sratch
COPY --from=builder /bin/example /bin/example
CMD [ "/bin/example" ]
```

```bash
$ docker build -t app:test . -f Dockerfile5
$ docker run app:test --name app
```

