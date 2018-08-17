package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func serveBackend(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	})
	// backend listen on port 9999
	fmt.Println("backend:", http.ListenAndServe(":9999", mux))
}

func main() {
	addr := "http://localhost:9999"
	rpURL, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	go serveBackend(addr)

	mux := http.NewServeMux()
	mux.Handle("/", httputil.NewSingleHostReverseProxy(rpURL)) // simple reverse

	// proxy exposes port 80
	fmt.Println("proxy:", http.ListenAndServe(":80", mux)) // needs sudo/root permission
}
