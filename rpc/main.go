package main

import (
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/ShevaXu/playground/rpc/proto"
)

var addr = ":9090" // fixed so far

func main() {
	// server
	handler := rpc.NewServer()
	rcvr := proto.NewCounter()

	if err := handler.Register(rcvr); err != nil {
		log.Fatal("Register handler", err)
	}

	handler.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	log.Println("Listen on addr", addr)
	go func() {
		log.Fatal(http.ListenAndServe(addr, handler))
	}()

	time.Sleep(time.Second)

	// client
	cl, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("Fail to dial", addr, err)
	}

	var ret int
	if err := cl.Call("Counter.Add", 2, &ret); err != nil {
		log.Println("RPC error", err)
	} else {
		log.Println("RPC result", ret)
	}
	if err = cl.Call("Counter.Get", proto.NoArgs{}, &ret); err != nil {
		log.Println("RPC error", err)
	} else {
		log.Println("RPC result", ret)
	}
	if err = cl.Call("Counter.Clear", proto.NoArgs{}, &proto.NoReturn{}); err != nil {
		log.Println("RPC error", err)
	}
	if err = cl.Call("Counter.Get", proto.NoArgs{}, &ret); err != nil {
		log.Println("RPC error", err)
	} else {
		log.Println("RPC result", ret)
	}
}
