# proxy

A simple proxy server for HTTP(S).

## Test

```sh
go run main.go

curl -Lv -x http://localhost:8082 https://baidu.com
```

## HTTP Tunnel

From [wikipedia](https://en.wikipedia.org/wiki/HTTP_tunnel)

> A variation of HTTP tunneling when behind an HTTP proxy server is to use the ["CONNECT"](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT) HTTP method. In this mechanism, the client asks an HTTP proxy server to forward the TCP connection to the desired destination. The server then proceeds to make the connection on behalf of the client. Once the connection has been established by the server, the proxy server continues to proxy the TCP stream to and from the client.

## Client

Go 1.10 supports HTTPs proxy [*now*](https://medium.com/@mlowicki/https-proxies-support-in-go-1-10-b956fb501d6b). Check `client/main.go`. (However, test with Go 1.9.3 works as well.)

## What's More

#### Hop-by-hop headers

https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#hbh

- Proxy-*
- Upgrade
- Keep-Alive
- Transfer-Encoding

#### Forward headers

- X-Forwarded-Proto
- X-Forwarded-Host

#### Streaming

## Refs

Inspired by the [blog](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c).

### More libs

General purpose

- https://github.com/elazarl/goproxy
- https://github.com/ginuerzh/gost
- https://github.com/jpillora/chisel

Reverse proxy

- https://github.com/containous/traefik
- https://github.com/fatedier/frp
- https://github.com/vulcand/oxy

For bypass network restriction

- https://github.com/getlantern/lantern
- https://github.com/v2ray/v2ray-core
- https://github.com/cyfdecyf/cow
- https://github.com/yinghuocho/firefly-proxy
