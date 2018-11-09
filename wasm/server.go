package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port = flag.Int("p", 9080, "listen port")
	dir  = flag.String("dir", ".", "directory to serve")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	log.Printf("listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(*dir))))
}
