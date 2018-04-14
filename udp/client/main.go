package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//simple write
	conn.Write([]byte("Hello from client"))

	//simple Read
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	log.Println(string(buffer[:n]))
}
