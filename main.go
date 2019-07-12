package main

import (
	"fmt"
	"github.com/aeidelos/go-sse-notification/client"
	"github.com/aeidelos/go-sse-notification/server"
	"log"
	"net/http"
)

func main() {
	broker := server.NewServer()
	publisher := server.NewEvent(broker)
	http.Handle("/event", http.HandlerFunc(publisher.PublishEventHTTP))
	http.Handle("/listen", broker)
	http.Handle("/", http.HandlerFunc(client.DisplayWebPage))
	fmt.Println("starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
