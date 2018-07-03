# [Twirp](https://twitchtv.github.io/twirp/)

An alternative to gRPC.

## Example

Follow https://twitchtv.github.io/twirp/docs/example.html.

```sh
# generate code
$ protoc --proto_path=$GOPATH/src:. --twirp_out=. --go_out=. ./rpc/haberdasher/service.proto
# serve :8080
go run server.go
# another shell
curl -X POST -H "Content-Type:application/json" -d '{"inches": 10}' -v http://localhost:8080/twirp/twirp.example.haberdasher.Haberdasher/MakeHat
```
