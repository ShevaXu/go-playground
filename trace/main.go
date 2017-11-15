package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"time"
)

// RequestInfo groups the request hook info for print.
type RequestInfo struct {
	DNSStart httptrace.DNSStartInfo
	DNSDone  httptrace.DNSDoneInfo
	GotConn  httptrace.GotConnInfo
	Net      string
	Addr     string
}

// from github.com/davecheney/httpstat
const (
	HTTPSTemplate = "\n==========\n" +
		`  DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer` + "\n" +
		`[%s  |     %s  |    %s  |        %s  |       %s  ]` + "\n" +
		`            |                |               |                   |                  |` + "\n" +
		`   name-lookup:%s      |               |                   |                  |` + "\n" +
		`                       connect:%s     |                   |                  |` + "\n" +
		`                                   pre-transfer:%s         |                  |` + "\n" +
		`                                                     start-transfer:%s        |` + "\n" +
		`                                                                                total:%s` +
		"\n==========\n\n"

	HTTPTemplate = "\n==========\n" +
		`   DNS Lookup   TCP Connection   Server Processing   Content Transfer` + "\n" +
		`[ %s  |     %s  |        %s  |       %s  ]` + "\n" +
		`             |                |                   |                  |` + "\n" +
		`    name-lookup:%s      |                   |                  |` + "\n" +
		`                        connect:%s         |                  |` + "\n" +
		`                                      start-transfer:%s        |` + "\n" +
		`                                                                 total:%s` +
		"\n==========\n\n"
)

// helper
func stringify(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Require url input")
	}

	uri, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("Invalid url", os.Args[1])
	}
	fmt.Println("Parsed URL", stringify(uri))

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	cl := http.Client{
		Transport: tr,
	}

	var (
		t0, t1, t2, t3, t4 time.Time
		info               RequestInfo
	)

	trace := httptrace.ClientTrace{
		DNSStart: func(i httptrace.DNSStartInfo) {
			t0 = time.Now()
			info.DNSStart = i
		},
		DNSDone: func(i httptrace.DNSDoneInfo) {
			t1 = time.Now()
			info.DNSDone = i
		},
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				// connecting to IP
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Fatalln("Unable to connect to host:", err)
			}
			t2 = time.Now()

			info.Net = net
			info.Addr = addr
		},
		GotConn: func(i httptrace.GotConnInfo) {
			t3 = time.Now()
			info.GotConn = i
		},
		GotFirstResponseByte: func() {
			t4 = time.Now()
		},
	}

	req, _ := http.NewRequest("GET", uri.String(), nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	resp, err := cl.Do(req)
	if err != nil {
		log.Fatalln("Do request error", err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Read body error", err)
	}

	t5 := time.Now()

	switch uri.Scheme {
	case "https":
		fmt.Printf(HTTPSTemplate,
			t1.Sub(t0), // dns lookup
			t2.Sub(t1), // tcp connection
			t3.Sub(t2), // tls handshake
			t4.Sub(t3), // server processing
			t5.Sub(t4), // content transfer
			t1.Sub(t0), // name-lookup
			t2.Sub(t0), // connect
			t3.Sub(t0), // pre-transfer
			t4.Sub(t0), // start-transfer
			t5.Sub(t0), // total
		)
	case "http":
		fmt.Printf(HTTPTemplate,
			t1.Sub(t0), // dns lookup
			t3.Sub(t1), // tcp connection
			t4.Sub(t3), // server processing
			t5.Sub(t4), // content transfer
			t1.Sub(t0), // name-lookup
			t3.Sub(t0), // connect
			t4.Sub(t0), // start-transfer
			t5.Sub(t0), // total
		)
	}

	fmt.Println(stringify(info))
}
