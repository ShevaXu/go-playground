# Cgo

> Cgo lets Go packages call C code.

More about [Cgo](https://golang.org/cmd/cgo/)

## Simple Example

Follow https://blog.golang.org/c-go-cgo

```sh
# The go tool recognizes the special "C" import and automatically uses cgo for those files.
go run basic/main.go
```

## More Advanced - Calling libavx from FFmpeg project.

```sh
cd ffmpeg
# put your own test.mp4 there
# run stream info probing
go run main.go
```

Requirements:

- pkg-config
- ffmpeg (v3.2.13 tested) installed (`configure; make; make install`)

For wrapper, consider using https://github.com/xlab/c-for-go or https://github.com/giorgisio/goav directly (need to tune for different FFmpeg versions).

## C++

Using C wrapper, or [SWIG](http://www.swig.org/Doc2.0/Go.html)

Examples: https://github.com/draffensperger/go-interlang
