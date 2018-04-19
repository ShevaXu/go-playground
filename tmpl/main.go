package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const simplePage = `
<html><body>
<h2>Hello {{$}}</h2>
</body></html>
`

var (
	port = flag.Int("port", 8089, "Service port")
	t    *template.Template
)

func simpleHandle(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "stranger"
	}
	if err := t.Execute(w, name); err != nil {
		log.Println(err)
	}
}

func main() {
	var err error
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	t, err = template.New("simple").Parse(simplePage)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/hello", simpleHandle)
	log.Println("Listening on http://localhost" + addr)
	log.Println(http.ListenAndServe(addr, nil))
}
