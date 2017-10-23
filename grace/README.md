# grace

Golang Server with Graceful Shutdown/Restart.

## Signals

`SIGTERM` will trigger a graceful shutdown;

`SIGINT` will trigger a fork/restart then shutdown.
(NOTE that `SIGINT` can only trigger once the fork for a process, following `SIGINT` will be ignored.)

## Usage

Just like you use `http.ListenAndServe()`:

```
err := grace.ListenAndServe(":9090", nil)
```

### Play with `/demo`

The basic handler sleeps w.r.t. the `duration` then returns:

```
$ curl http://127.0.0.1:9090?duration=1s
started at 2017-02-21 23:26:19.135563128 +0800 CST slept for 1.000 seconds from pid 59715.
```

So try `curl http://127.0.0.1:9090?duration=5s` and immediately sends `SIGINT`: 

```
$ kill -2 $pid
2017/02/21 23:27:59 Server 59742 forked
2017/02/21 23:27:59 Serving :9090 with pid 59747.
2017/02/21 23:28:02 Server 59742 shutdown
2017/02/21 23:28:02 Server 59742 stoped. Error - accept tcp [::]:9090: use of closed network connection

// curl should still returns
started at 2017-02-21 23:28:02.079373063 +0800 CST slept for 5.000 seconds from pid 59742.
```

### Go 1.8

In `/go1.8`, grace is re-implemented with Go 1.8 with `Server.Shutdown()` added.

## TODOs

* (unit) tests
* configurable (keepAlivePeriod, signals)
* features (hammerTime, timeouts, signal hooks ...)
* documentation

## References

* The blog post http://grisha.org/blog/2014/06/03/graceful-restart-in-golang

### Other Implementations

* endless https://github.com/fvbock/endless/
* gracehttp https://github.com/tabalt/gracehttp
* https://github.com/facebookgo/grace
* graceful https://github.com/tylerb/graceful
