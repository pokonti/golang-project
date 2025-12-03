package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ClientManager keeps track of all the active websocket connections
type ClientManager struct {
	// A map of UserID -> WebSocket Connection
	clients map[int64]*websocket.Conn
	// Mutex to protect the map from concurrent writes
	lock sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[int64]*websocket.Conn),
	}
}

// AddClient registers a user connection
func (manager *ClientManager) AddClient(userID int64, conn *websocket.Conn) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.clients[userID] = conn
}

// RemoveClient unregisters a user connection
func (manager *ClientManager) RemoveClient(userID int64) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if conn, ok := manager.clients[userID]; ok {
		conn.Close()
		delete(manager.clients, userID)
	}
}

// GetClient returns the connection for a specific user
func (manager *ClientManager) GetClient(userID int64) (*websocket.Conn, bool) {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	conn, ok := manager.clients[userID]
	return conn, ok
}
