// Package gosocketify is a package with server side implementation for socketify.
// It provides a router for handling WebSocket connections and HTTP requests.
package gosocketify

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Router struct handles upgrading a websocket connection with all CORS allowed.
// It manages incoming WebSocket connections and integrates with a connection manager
// to keep track of active connections.
type Router struct {
	connectionManager *ConnectionManager
}

// NewRouter creates a new instance of the Router.
// It initializes the Router with the provided connection manager,
// sets up the appropriate routes for WebSocket and HTTP,
// and returns an http.Handler that can be used to handle incoming requests.
func NewRouter(connectionManager *ConnectionManager) http.Handler {
	rm := &Router{
		connectionManager: connectionManager,
	}

	r := mux.NewRouter()

	// Define routes for WebSocket
	r.HandleFunc("/ws", rm.handleWebSocket)

	// Define HTTP routes
	r.HandleFunc("/", rm.handleHome).Methods("GET")

	http.Handle("/", r)

	return r
}

// handleWebSocket upgrades the connection from HTTP to WebSocket.
// It accepts incoming WebSocket connections, applies CORS checks,
// and adds the connection to the connection manager.
// If the upgrade process fails, it responds with an error.
func (rm *Router) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}

	rm.connectionManager.AddConnection(conn)
}

// handleHome responds with a welcome message to HTTP GET requests at the root path.
// This can serve as a simple health check endpoint for the WebSocket server.
func (rm *Router) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the WebSocket server!"))
}
