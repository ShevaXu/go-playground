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
protoc -I contacts/ contacts/contacts.proto --go_out=plugins=grpc:contacts
```

Run the server & client:

```sh
go run server.go

go run client/main.go
```

## Refs

### Proto3

- Fields are 'optional' by default, "required/optional" are removed ([why](https://github.com/google/protobuf/issues/2497)).
- Explicit default values are not allowed.
