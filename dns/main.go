package main

import (
	"flag"
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

var res = flag.String("r", "1.1.1.1", "DNS resolver to forward the query")

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	resolver := net.UDPAddr{IP: net.ParseIP(*res), Port: 53}
	log.Printf("DNS proxy @%s\n=====>\n", resolver.String())

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	panicErr(err)
	defer conn.Close()

	var asker *net.UDPAddr
	for {
		// read
		buf := make([]byte, 512)
		_, addr, err := conn.ReadFromUDP(buf)
		panicErr(err)
		// parse
		var m dnsmessage.Message
		panicErr(m.Unpack(buf))
		log.Printf("Got udp: %+v\n", m)
		// ignore invalid query
		if len(m.Questions) == 0 {
			continue
		}
		// real magic
		if !m.Header.Response { // asking
			asker = addr // memorise; TODO: it's not for concurrent use
			log.Println("Forward to", resolver)
			_, err = conn.WriteToUDP(buf, &resolver)
			panicErr(err)
		} else { // answering
			log.Println("Answer", *asker)
			_, err = conn.WriteToUDP(buf, asker)
			panicErr(err)
			log.Println("Session done\n=====>")
		}
	}
}
