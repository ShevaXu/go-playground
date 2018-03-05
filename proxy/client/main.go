package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	port = flag.Int("port", 8082, "Proxy port")
	addr = flag.String("addr", "https://baidu.com", "Request address")
)

func main() {
	flag.Parse()

	u, err := url.Parse(fmt.Sprintf("http://localhost:%d", *port))
	if err != nil {
		panic(err)
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(*addr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(dump))
}
