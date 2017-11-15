# HTTP Tracing with `net/http/httptrace`

The Go [blog](https://blog.golang.org/http-tracing) introduce `httptrace` (hooks) for these HTTP events:

- Connection creation
- Connection reuse
- DNS lookups
- Writing the request to the wire
- Reading the response

## Usage

```shell
$ go run main.go https://baidu.com/
```

## Refs

> Imitation is the sincerest form of flattery.

https://github.com/davecheney/httpstat
