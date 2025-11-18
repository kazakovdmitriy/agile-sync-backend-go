package websocket

import (
	"backend_go/internal/infrastructure/config"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager управляет только WebSocket соединениями
type ConnectionManager struct {
	// sessionID -> []*websocket.Conn
	sessions map[string][]*websocket.Conn
	mutex    sync.RWMutex
	upgrader websocket.Upgrader
}

// NewConnectionManager создает новый менеджер
func NewConnectionManager(cfg *config.Config) *ConnectionManager {
	return &ConnectionManager{
		sessions: make(map[string][]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin:       checkOrigin(cfg),
			HandshakeTimeout:  cfg.WebSocket.HandshakeTimeout,
			ReadBufferSize:    cfg.WebSocket.ReadBufferSize,
			WriteBufferSize:   cfg.WebSocket.WriteBufferSize,
			EnableCompression: cfg.WebSocket.EnableCompression,
		},
	}
}

// Connect подключает соединение к сессии
func (cm *ConnectionManager) Connect(sessionID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.sessions[sessionID]; !exists {
		cm.sessions[sessionID] = []*websocket.Conn{}
	}

	cm.sessions[sessionID] = append(cm.sessions[sessionID], conn)
	log.Printf("Client connected to session %s. Total connections: %d", sessionID, len(cm.sessions[sessionID]))
}

// Disconnect отключает соединение от сессии
func (cm *ConnectionManager) Disconnect(sessionID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connections, exists := cm.sessions[sessionID]
	if !exists {
		return
	}

	for i, connection := range connections {
		if connection == conn {
			cm.sessions[sessionID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	if len(cm.sessions[sessionID]) == 0 {
		delete(cm.sessions, sessionID)
	}

	log.Printf("Client disconnected from session %s. Remaining connections: %d", sessionID, len(cm.sessions[sessionID]))
}

// Broadcast рассылает сообщение всем соединениям в сессии
func (cm *ConnectionManager) Broadcast(sessionID string, message interface{}) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connections, exists := cm.sessions[sessionID]
	if !exists {
		return nil
	}

	var errors []error
	for _, conn := range connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Broadcast error: %v", err)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// SendTo отправляет сообщение конкретному соединению
func (cm *ConnectionManager) SendTo(conn *websocket.Conn, message interface{}) error {
	return conn.WriteJSON(message)
}

// GetUpgrader возвращает upgrader для использования в хендлере
func (cm *ConnectionManager) GetUpgrader() websocket.Upgrader {
	return cm.upgrader
}

// GetSessionCount возвращает количество активных сессий
func (cm *ConnectionManager) GetSessionCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.sessions)
}

// GetConnectionCount возвращает общее количество соединений
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	total := 0
	for _, connections := range cm.sessions {
		total += len(connections)
	}
	return total
}
