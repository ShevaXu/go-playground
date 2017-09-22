# Server-Sent Event Demo

> Server-sent events ([SSE](https://en.wikipedia.org/wiki/Server-sent_events)) is a technology where a browser receives automatic updates from a server via HTTP connection.

This go implementation is inspired by the [gist](https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c), its [blog](https://robots.thoughtbot.com/writing-a-server-sent-events-server-in-go) and the [fixes](https://gist.github.com/schmohlio/d7bdb255ba61d3f5e51a512a7c0d6a85).

## Added Features

- Sub-route at `/sse`
- Provided `index.html`
- Broadcast messages in goroutine for each client

## Web SSE

[Using server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)

```javascript
var es = new EventSource("http://localhost:3000/sse")

es.onmessage = function(msg) {
    console.log(msg)
}
```
