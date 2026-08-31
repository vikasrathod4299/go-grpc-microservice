// Package hub provides an in-memory hub for managing active WebSocket connections indexed by user_id (Driver or Rider).
package hub

/*
================================================================================
FILE: internal/gateway/hub/hub.go
================================================================================

PURPOSE:
Maintains an in-memory dictionary of active WebSocket connections indexed by user_id (Driver or Rider).
Allows backend event processors to send targeted real-time push messages to specific users.

LEARNING GO CONCEPTS:
- Concurrent map management with Mutex (`sync.RWMutex`) or channel select loop.
- Go channel event loop pattern (`register`, `unregister`, `broadcast`).

WHAT YOU NEED TO IMPLEMENT HERE:
1. `Client` struct representing an active connection:
   - UserID string
   - Role string ("driver" or "rider")
   - Conn *websocket.Conn
   - Send chan []byte
2. `Hub` struct:
   - clients map[string]*Client // key is UserID
   - register chan *Client
   - unregister chan *Client
   - broadcast chan Message
3. `Run()` method:
   - Infinite loop with `select` handling register, unregister, and broadcast events.
4. `SendToUser(userID string, payload []byte) bool`:
   - Safely pushes a message to a specific user's send channel.
================================================================================
*/

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Role   string // "driver" or "rider"
	Conn   *websocket.Conn
	Send   chan []byte
}

// Message represents a targeted payload to a specific user.
type Message struct {
	UserID  string
	Payload []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

// NewHub initializes and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

// Run starts the infinite event loop for the Hub.
// This must be run as a goroutine: `go hub.Run()`
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Add the new client to the map
			h.clients[client.UserID] = client

		case client := <-h.unregister:
			// Remove the client and safely close their channel
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}

		case msg := <-h.broadcast:
			// Look up the specific user and send them the payload
			if client, ok := h.clients[msg.UserID]; ok {
				select {
				case client.Send <- msg.Payload:
					// Message sent successfully to the client's internal channel
				default:
					// If the client's send channel is blocked/full, kick them
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
		}
	}
}

// SendToUser provides a safe, public method for other packages to push messages.
// It wraps the payload in a Message struct and sends it to the central broadcast channel.
func (h *Hub) SendToUser(userID string, payload []byte) bool {
	h.broadcast <- Message{
		UserID:  userID,
		Payload: payload,
	}
	return true
}
