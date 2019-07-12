package server

import (
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	Notifier       chan string
	IncomingClient chan chan string
	ExitingClient  chan chan string
	ActiveClient   map[chan string]bool
}

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan string)
	broker.IncomingClient <- messageChan
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		broker.ExitingClient <- messageChan
		log.Println("http connection closed")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {
		msg, open := <-messageChan
		if !open {
			break
		}
		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}
	log.Println("finish http request ", r.URL.Path)

}

func (broker *Broker) Listen() {
	for {
		select {
		case s := <-broker.IncomingClient:
			broker.ActiveClient[s] = true
			log.Printf("add client. %d registered clients", len(broker.ActiveClient))
		case s := <-broker.ExitingClient:
			delete(broker.ActiveClient, s)
			close(s)
			log.Printf("remove client. %d registered clients", len(broker.ActiveClient))
		case event := <-broker.Notifier:
			for client, _ := range broker.ActiveClient {
				client <- event
			}
		}
	}
}

func NewServer() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan string),
		IncomingClient: make(chan chan string),
		ExitingClient:  make(chan chan string),
		ActiveClient:   make(map[chan string]bool),
	}
	go broker.Listen()
	return broker
}
