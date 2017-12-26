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

### Makefile

*Ref.:* https://blog.gopheracademy.com/advent-2017/make/

## Beyond Basics

*Ref.:* https://blog.gopheracademy.com/advent-2017/go-grpc-beyond-basics/

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

### More

- Tracing with `golang.org/x/net/trace` package (`grpc.EnableTracing = true`)
- [Backoff](https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md)
- [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [gRPC Web](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md)

## Details

- [godoc](https://godoc.org/google.golang.org/grpc)
- [Proto3](https://developers.google.com/protocol-buffers/docs/proto3)

### Proto3 Design

- All fields are "optional" by design - "required/optional" keywords are removed.

> We dropped required fields in proto3 because required fields are generally considered harmful and violating protobuf's compatibility semantics.
(https://github.com/google/protobuf/issues/2497)).

- Explicit default values are not allowed. See [default values](https://developers.google.com/protocol-buffers/docs/proto3#default).

> ... primitive (non-message) fields are no longer nullable. It's better/more accurate to not even think of "unset"; i.e., a new message object already has every primitive field set. The intent, as I understand it, was to make a message more like a plain struct. https://github.com/google/protobuf/issues/359#issuecomment-101756694

- Use wrapper ([well-known types](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf)) to test presence.

See https://github.com/google/protobuf/blob/master/src/google/protobuf/wrappers.proto

```protobuf
// Wrappers for primitive (non-message) types. These types are useful
// for embedding primitives in the `google.protobuf.Any` type and for places
// where we need to distinguish between the absence of a primitive
// typed field and its default value.
```

### gPRC (Proto3) in Go

https://developers.google.com/protocol-buffers/docs/reference/go-generated

```protobuf
message Bar {}

message Baz {
  Bar foo = 1;
}

message Baz2 {
  repeated Bar foo = 1;
}
```

The compiler will generate a Go struct

```go
type Baz struct {
  Foo *Bar
}

type Baz2 struct {
  Foo []*Bar
}
```

**Thus message fields will be set to `nil` if unset.**
