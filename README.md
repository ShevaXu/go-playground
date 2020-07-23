# Go Playground

Explorations of new Go features and Go implementations of web tech.

## Updates

#### grace

- Go 1.9: `func (srv *Server) RegisterOnShutdown(f func())`
- Go 1.13

```go
type Server struct {
    // ...

    // BaseContext optionally specifies a function that returns
    // the base context for incoming requests on this server.
    // The provided Listener is the specific Listener that's
    // about to start accepting requests.
    // If BaseContext is nil, the default is context.Background().
    // If non-nil, it must return a non-nil context.
    BaseContext func(net.Listener) context.Context // Go 1.13

    // ConnContext optionally specifies a function that modifies
    // the context used for a new connection c. The provided ctx
    // is derived from the base context and has a ServerContextKey
    // value.
    ConnContext func(ctx context.Context, c net.Conn) context.Context // Go 1.13
}
```
