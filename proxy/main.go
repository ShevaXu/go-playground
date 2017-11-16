package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8082", "Listen address")

func tunnel(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		dst, src net.Conn
		hijacker http.Hijacker
		ok       bool
	)

	dst, err = net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println("Dial succeed")

	if hijacker, ok = w.(http.Hijacker); !ok {
		http.Error(w, "Hijacking is not supported", http.StatusInternalServerError)
		return
	}

	if src, _, err = hijacker.Hijack(); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	log.Println("Hijack succeed")

	// Reply OK to CONNECT request
	src.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	log.Println("Start tunneling")

	// bi-directional
	// TODO: error handling
	go io.Copy(dst, src)
	io.Copy(src, dst)

	log.Println("Finish tunneling")

	return
}

func handleHTTP(w http.ResponseWriter, req *http.Request) (err error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)
	written, err := io.Copy(w, resp.Body)
	log.Printf("Proxy http writes %d bytes", written)

	return
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func proxy(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == http.MethodConnect {
		err = tunnel(w, r)
	} else {
		err = handleHTTP(w, r)
	}
	if err != nil {
		log.Println("Error performing connect", err)
	}
}

func main() {
	flag.Parse()

	server := &http.Server{
		Addr:    *addr,
		Handler: http.HandlerFunc(proxy),
		// If TLSNextProto is not nil, HTTP/2 support is not enabled
		// automatically.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Println("Serving proxy from", *addr)
	log.Fatal(server.ListenAndServe())
}
