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
	ts   *template.Template
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

func headersHandle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	list := make([]string, 0, len(r.Header))
	for k, v := range r.Header {
		list = append(list, fmt.Sprintf("%s: %+v", k, v))
	}
	log.Println(list)
	ts.ExecuteTemplate(w, "header", nil)
	if err := ts.ExecuteTemplate(w, "main", list); err != nil {
		log.Println(err)
	}
	// another way
	footer := ts.Lookup("footer")
	footer.Execute(w, nil)
}

func main() {
	var err error
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	t, err = template.New("simple").Parse(simplePage)
	if err != nil {
		panic(err)
	}
	// multipart template from files
	ts, err = template.ParseGlob("tmpls/*.tmpl")
	// or
	// ts, err = template.ParseFiles("tmpls/header.tmpl", "tmpls/main.tmpl", "tmpls/footer.tmpl")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/hello", simpleHandle)
	http.HandleFunc("/headers", headersHandle)
	log.Println("Listening on http://localhost" + addr)
	log.Println(http.ListenAndServe(addr, nil))
}
