// Package gosocketify is a package with server side implementation for socketify.
// It includes message handling and event dispatching for WebSocket connections.
package gosocketify

import "log"

// Message struct represents a message sent over the WebSocket connection.
// It contains an event type and data associated with that event.
type Message struct {
	Event string      `json:"event"` // The type of event represented by this message.
	Data  interface{} `json:"data"`  // The data payload associated with the event.
}

// EventDispatcher interface defines methods for handling messages and events.
// It provides methods to dispatch messages to specific connections, broadcast
// messages to all clients, and send messages to specific connections.
type EventDispatcher interface {
	Dispatch(connId string, message Message) // Dispatch a message for a specific connection.
	Broadcast(message Message)               // Broadcast a message to all connected clients.
	Send(connectionID any, message Message)  // Send a message to a specific connection identified by its ID.
}

// ExampleEventDispatcher is a basic implementation of the EventDispatcher interface.
// It provides stub functionality for dispatching and handling events.
type ExampleEventDispatcher struct{}

// Dispatch processes incoming messages and handles them based on the event type.
// It contains logic to handle specific events such as "broadcast" and "ping".
func (ed *ExampleEventDispatcher) Dispatch(connId string, message Message) {
	switch message.Event {
	case "broadcast":
		// Placeholder for broadcasting a message to all clients.
		log.Println("Unimplemented broadcasting to all users")
	case "ping":
		// Placeholder for echoing a "pong" message back to the sender.
		log.Println("Unimplemented sending data back to user")

	default:
		// Log a message for unhandled event types.
		log.Printf("Event not handled: %s", message.Event)
	}
}
