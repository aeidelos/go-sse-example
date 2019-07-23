package server

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Event struct {
	broker *Broker
}

func NewEvent(broker *Broker) *Event {
	return &Event{broker: broker}
}

func (event *Event) PublishEventHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to parse http body", http.StatusInternalServerError)
		return
	}
	event.broker.Notifier <- requestBody
	log.Printf("retrieve event : " + string(requestBody))
	w.Write([]byte("event sent"))
}
