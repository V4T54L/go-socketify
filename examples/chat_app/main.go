package main

import (
	"log"
	"net/http"

	gosocketify "github.com/V4T54L/go-socketify"
)

func main() {
	connectionManager := gosocketify.NewConnectionManager(nil)
	eventDispatcher := newEventDispatcher(*connectionManager)
	connectionManager.EventDispatcher = eventDispatcher

	// Initialize the router
	router := gosocketify.NewRouter(connectionManager)

	// Start the HTTP and socketify server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

type eventDispatcher struct {
	connectionManager *gosocketify.ConnectionManager
}

func newEventDispatcher(cm gosocketify.ConnectionManager) *eventDispatcher {
	return &eventDispatcher{connectionManager: &cm}
}

func (e *eventDispatcher) Broadcast(msg gosocketify.Message) {
	e.connectionManager.Broadcast(msg)
}

func (e *eventDispatcher) Send(connectionID any, message gosocketify.Message) {
	id, ok := connectionID.(string)
	if !ok {
		log.Println("Failed to extract id string from connectionID")
		return
	}
	e.connectionManager.ToConnection(id, message)
}

func (e *eventDispatcher) Dispatch(connId string, message gosocketify.Message) {
	switch message.Event {
	case "broadcast":
		message.Event = "broadcastResponse"
		e.Broadcast(message)
	case "ping":
		message.Event = "pong"
		e.Send(connId, message)

	default:
		// Log a message for unhandled event types.
		log.Printf("Event not handled: %s", message.Event)
	}
}
