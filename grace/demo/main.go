package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ShevaXu/grace"
)

var pid int

func sleep(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration(r.FormValue("duration"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	time.Sleep(duration)

	fmt.Fprintf(
		w,
		"started at %s slept for %.3f seconds from pid %d.\n",
		time.Now(),
		duration.Seconds(),
		pid,
	)
}

func main() {
	pid = os.Getpid()
	addr := ":9090"
	log.Printf("Serving %s with pid %d.\n", addr, pid)

	http.HandleFunc("/", sleep) // defaultMux

	//srv := grace.GraceServer{Addr: addr}
	//err := srv.ListenAndServe()
	err := grace.ListenAndServe(addr, nil)

	log.Printf("Server %d stoped. Error - %v\n", pid, err)
}
