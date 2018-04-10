package main

import (
	"flag"
	"log"
	"net"
)

var (
	host = flag.String("h", "127.0.0.1", "udp host")
	port = flag.Int("p", 9090, "udp port")
)

func main() {
	flag.Parse()

	addr := net.UDPAddr{
		IP:   net.ParseIP(*host),
		Port: *port,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		//simple read
		buffer := make([]byte, 1024)
		n, from, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		log.Println("Read", n, "bytes:", string(buffer))

		//simple write
		_, err = conn.WriteTo([]byte("Reply from server\n"), from)
		if err != nil {
			panic(err)
		}
	}
}
