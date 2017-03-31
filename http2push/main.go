package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const indexHTML = `<html>
<head>
	<title>Hello Gopher</title>
	<script src="/static/app.js"></script>
	<link rel="stylesheet" href="/static/style.css"">
</head>
<body>
Hello gopher from HTTP2 Push!
</body>
</html>
`

var addr = flag.String("addr", ":8081", "Listen address")

func HelloPush(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		// serve only "/"
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}
		if err := pusher.Push("/static/app.js", options); err != nil {
			log.Printf("Failed to push js: %v", err)
		}
		if err := pusher.Push("/static/style.css", options); err != nil {
			log.Printf("Failed to push css: %v", err)
		}
	}
	fmt.Fprint(w, indexHTML)
}

const staticPrefix = "/static/"

func main() {
	flag.Parse()

	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HelloPush)

	log.Println("Serving HTTP2 from", *addr)
	log.Fatal(http.ListenAndServeTLS(*addr, "cert.pem", "key.pem", nil))
}
