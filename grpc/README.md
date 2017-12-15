# [gRPC](https://grpc.io/)

> In gRPC a client application can directly call methods on a server application on a different machine as if it was a local object, making it easier for you to create distributed applications and services.

![](https://grpc.io/img/landing-2.svg)

## Requirements

- `protoc` complier for [protobuf](https://github.com/golang/protobuf)
- gRPC: `go get -u google.golang.org/grpc`
- *protoc-gen-go* plugin: `go get -u github.com/golang/protobuf/protoc-gen-go`

## Demo

Build the stubs:

```sh
make compile
# protoc -I contacts/ contacts/contacts.proto --go_out=plugins=grpc:contacts
```

Run the server & client:

```sh
go run server.go

go run client/main.go
```

## Beyond Basics

### Secure gRPC

Generate the certificate:

```sh
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
# IMPORTANT: MUST specify server name after the prompt
# Common Name (e.g. server FQDN or YOUR name) []:localhost
```

Add [credential](https://godoc.org/google.golang.org/grpc/credentials) options:

```go
// server
cred, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
s := grpc.NewServer(grpc.Creds(cred))

// client
cred, err := credentials.NewClientTLSFromFile("cert.pem", "")
conn, err := grpc.Dial("localhost:15001", grpc.WithTransportCredentials(cred))
```

### Interceptor (Middleware)

Interceptor for [Unary RPC](https://grpc.io/docs/guides/concepts.html#unary-rpc) (single request) with [metadata](https://godoc.org/google.golang.org/grpc/metadata) through `context.Context`.



## Refs

- https://godoc.org/google.golang.org/grpc
- https://blog.gopheracademy.com/advent-2017/go-grpc-beyond-basics/

### Proto3

https://developers.google.com/protocol-buffers/docs/proto3

- Fields are 'optional' by default, "required/optional" are removed ([why](https://github.com/google/protobuf/issues/2497)).
- Explicit default values are not allowed.

### Makefile

Ref to https://blog.gopheracademy.com/advent-2017/make/
