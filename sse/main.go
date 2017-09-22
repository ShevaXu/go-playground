// v2 of the great example of SSE in go by @ismasan.
// includes fixes:
//    * infinite loop ending in panic
//    * closing a client twice
//    * potentially blocked listen() from closing a connection during multiplex step.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range clients` started.
const patience time.Duration = time.Second * 1

// Example SSE server in Golang.
//     $ go run sse.go

type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func NewBroker() (b *Broker) {
	// Instantiate a broker
	b = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go b.listen()

	return
}

func (b *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Make sure that the writer supports flushing
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	b.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		b.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-notify:
			return
		default:
			// Write to the ResponseWriter
			// Server Sent Events compatible
			fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

			// Flush the data immediately instead of buffering it for later
			flusher.Flush()
		}
	}
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			// A new client has connected.
			// Register their message channel
			b.clients[s] = true
			log.Printf("Client added. %d registered clients", len(b.clients))
		case s := <-b.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(b.clients, s)
			log.Printf("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range b.clients {
				select {
				case clientMessageChan <- event:
				case <-time.After(patience):
					log.Print("Skipping client.")
				}
			}
		}
	}
}

func main() {
	//
	b := NewBroker()

	mux := http.NewServeMux()

	mux.Handle("/sse", b)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			b.Notifier <- []byte(eventString)
		}
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe(":3000", mux))
}