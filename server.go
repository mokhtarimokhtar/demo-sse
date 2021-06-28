package main

import (
	"fmt"
	"log"
	"net/http"
	"sse/config"
	"time"
)

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

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
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
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	// notify := rw.(http.CloseNotifier).CloseNotify()
	notify := req.Context().Done()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {
		// Write to the ResponseWriter
		// Server Sent Events compatible
		_, _ = fmt.Fprintf(rw, "%s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}

func main() {
	conf := config.Configuration{}
	err := conf.ReadConfig("conf.json")
	if err != nil {
		log.Println(err)
	}
	clock := time.Second * time.Duration(conf.Clock.Refresh)

	broker := NewServer()
	index := 0
	go func() {
		for {
			index += 1
			time.Sleep(clock)

			now := time.Now().UTC().Format(time.RFC3339)
			evt := "event: time"
			data := fmt.Sprintf("data: %v", now)
			//id := fmt.Sprintf("id:%v", index)
			//eventString := fmt.Sprintf("%v\n%v\n%v\n",data, evt, id)
			eventString := fmt.Sprintf("%v\n%v\n", evt, data)

			broker.Notifier <- []byte(eventString)
		}
	}()

	http.HandleFunc("/clocktimes", broker.ServeHTTP)

	fmt.Println("Starting Server Sent Event: ", conf.Server.Host + ":" + conf.Server.Port)
	fmt.Printf("Clock Refresh: %v sec\n", conf.Clock.Refresh)
	log.Fatal("HTTP server error: ", http.ListenAndServe(conf.Server.Host + ":" + conf.Server.Port, nil))
	//log.Fatal("HTTP server error: ", http.ListenAndServe("0.0.0.0:3000", nil))
}
