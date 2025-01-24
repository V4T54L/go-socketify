// Package gosocketify is a package with server side implementation for socketify.
// It includes connection management and event dispatching for WebSocket clients.
package gosocketify

import (
	"log"

	"github.com/gorilla/websocket"
)

// Connection struct represents a WebSocket connection.
// It includes an ID to identify the connection, the actual WebSocket connection,
// and a reference to the ConnectionManager responsible for managing this connection.
type Connection struct {
	ID      string
	Socket  *websocket.Conn
	Manager *ConnectionManager
}

// ConnectionManager struct manages all active WebSocket connections.
// It keeps track of connections and facilitates broadcasting and emitting messages
// to connected clients.
type ConnectionManager struct {
	connections     map[*Connection]bool
	EventDispatcher EventDispatcher
}

// NewConnectionManager creates a new instance of ConnectionManager.
// It initializes an empty connections map and sets the event dispatcher used for
// handling incoming messages from the connections.
func NewConnectionManager(eventDispatcher EventDispatcher) *ConnectionManager {
	connManager := &ConnectionManager{
		connections:     make(map[*Connection]bool),
		EventDispatcher: eventDispatcher,
	}

	return connManager
}

// AddConnection adds a new WebSocket connection to the ConnectionManager.
// It creates a new Connection instance, assigns it a unique ID, and starts
// handling messages for that connection in a separate goroutine.
// Returns the newly created Connection.
func (cm *ConnectionManager) AddConnection(conn *websocket.Conn) *Connection {
	connection := &Connection{
		ID:      generateConnectionID(),
		Socket:  conn,
		Manager: cm,
	}

	cm.connections[connection] = true

	go cm.handleMessages(connection)

	return connection
}

// Broadcast sends a message to all currently connected clients.
// It iterates through all active connections and attempts to send the message.
// Errors encountered during sending are logged but do not interrupt the broadcast process.
func (cm *ConnectionManager) Broadcast(message Message) {
	for conn := range cm.connections {
		err := conn.Socket.WriteJSON(message)
		if err != nil {
			log.Printf("Error broadcasting message to %s: %v", conn.ID, err)
			continue
		}
	}
}

// Emit sends a message to a specific connection identified by the Connection instance.
// If an error occurs while sending the message, it logs the error.
func (cm *ConnectionManager) Emit(conn *Connection, message Message) {
	err := conn.Socket.WriteJSON(message)
	if err != nil {
		log.Printf("Error sending message to %s: %v", conn.ID, err)
	}
}

// ToConnection sends a message to a specific connection using its ID.
// It searches the active connections for the matching ID and emits the message
// to the corresponding connection.
func (cm *ConnectionManager) ToConnection(id string, message Message) {
	var conn Connection
	for c := range cm.connections {
		if c.ID == id {
			conn = *c
			break
		}
	}
	cm.Emit(&conn, message)
}

// handleMessages listens for incoming messages on the given connection.
// It reads messages in a loop, dispatches them using the EventDispatcher,
// and removes the connection when done or in case of an error.
func (cm *ConnectionManager) handleMessages(conn *Connection) {
	defer func() {
		cm.RemoveConnection(conn)
	}()

	for {
		var message Message
		err := conn.Socket.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Handle the message and dispatch events
		cm.EventDispatcher.Dispatch(conn.ID, message)
	}
}

// RemoveConnection removes the specified connection from the ConnectionManager.
// It closes the WebSocket connection and cleans up the underlying resources.
func (cm *ConnectionManager) RemoveConnection(conn *Connection) {
	delete(cm.connections, conn)
	conn.Socket.Close()
}
