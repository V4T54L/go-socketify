// Package gosocketify is a package with server side implementation for socketify.
// It provides utilities and middlewares to handle connections and messages.
package gosocketify

import (
	"log"
)

// Middleware type definition.
// A function that takes a connection and a message as input and returns a new message.
// This function is meant to be used as a pipeline to transform or modify the messages
// before they are processed by the server.
type MiddlewareFunc func(Connection, Message) Message

// LoggingMiddleware logs data received and the connection id.
//
// This middleware is useful for debugging and monitoring purposes.
// It prints the connection id and the received message to the console.
//
// The message returned is the original message, so you can use this middleware
// in any position in your pipeline without affecting the original message.
func LoggingMiddleware(conn Connection, message Message) Message {
	log.Printf("Received message from %s: %v", conn.ID, message)
	return message
}
