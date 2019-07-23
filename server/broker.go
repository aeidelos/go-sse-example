package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Broker struct {
	Notifier       chan []byte
	IncomingClient chan chan []byte
	ExitingClient  chan chan []byte
	ActiveClient   map[chan []byte]bool
}

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan []byte)
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
		var buffer bytes.Buffer

		message := string(msg)
		if len(message) > 0 {
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Replace(message, "\n", "\ndata: ", -1)))
		}
		buffer.WriteString("\n")
		fmt.Fprintf(w, "data: %s\n\n", buffer.String())
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
		Notifier:       make(chan []byte, 1),
		IncomingClient: make(chan chan []byte),
		ExitingClient:  make(chan chan []byte),
		ActiveClient:   make(map[chan []byte]bool),
	}
	go broker.Listen()
	return broker
}
