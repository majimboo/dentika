package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// WebSocket connection manager
type WSManager struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
}

var wsManager = &WSManager{
	clients: make(map[*websocket.Conn]bool),
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Token   string      `json:"token,omitempty"`
}

// Notification message structure
type NotificationMessage struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	UserID      uint      `json:"user_id,omitempty"`
	ClinicID    uint      `json:"clinic_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	IsRead      bool      `json:"is_read"`
	Priority    string    `json:"priority,omitempty"`
	ActionURL   string    `json:"action_url,omitempty"`
	ActionLabel string    `json:"action_label,omitempty"`
}

// Add client to manager
func (wm *WSManager) addClient(conn *websocket.Conn) {
	wm.clientsMu.Lock()
	defer wm.clientsMu.Unlock()
	wm.clients[conn] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(wm.clients))
}

// Remove client from manager
func (wm *WSManager) removeClient(conn *websocket.Conn) {
	wm.clientsMu.Lock()
	defer wm.clientsMu.Unlock()
	delete(wm.clients, conn)
	log.Printf("WebSocket client disconnected. Total clients: %d", len(wm.clients))
}

// Broadcast message to all connected clients
func (wm *WSManager) broadcast(message WSMessage) {
	wm.clientsMu.RLock()
	defer wm.clientsMu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}

	for conn := range wm.clients {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Failed to send message to client: %v", err)
			// Remove broken connection
			go wm.removeClient(conn)
		}
	}
}

// Broadcast to specific user (if we had user tracking)
func (wm *WSManager) broadcastToUser(userID uint, message WSMessage) {
	// For now, broadcast to all clients
	// In a production app, you'd track which connections belong to which users
	wm.broadcast(message)
}

// WebSocket handler
func WebSocketHandler(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		// Add client to manager
		wsManager.addClient(conn)
		defer wsManager.removeClient(conn)

		// Handle incoming messages
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			if messageType == websocket.TextMessage {
				var wsMsg WSMessage
				if err := json.Unmarshal(message, &wsMsg); err != nil {
					log.Printf("Failed to parse WebSocket message: %v", err)
					continue
				}

				// Handle different message types
				switch wsMsg.Type {
				case "auth":
					handleAuth(conn, wsMsg)
				case "ping":
					handlePing(conn)
				default:
					log.Printf("Unknown message type: %s", wsMsg.Type)
				}
			}
		}
	})(c)
}

// Handle authentication
func handleAuth(conn *websocket.Conn, msg WSMessage) {
	if msg.Token == "" {
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "No token provided",
		})
		return
	}

	// Parse and validate JWT token
	token, err := jwt.Parse(msg.Token, func(token *jwt.Token) (interface{}, error) {
		// Use your JWT secret key here
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "Invalid token",
		})
		return
	}

	// Send success response
	conn.WriteJSON(fiber.Map{
		"type":    "auth_success",
		"message": "Authenticated successfully",
	})
}

// Handle ping messages
func handlePing(conn *websocket.Conn) {
	conn.WriteJSON(fiber.Map{
		"type":      "pong",
		"timestamp": time.Now().Unix(),
	})
}

// Public functions for sending notifications

// Send appointment reminder notification
func SendAppointmentReminder(appointmentID uint, patientName string, minutesUntil int) {
	message := WSMessage{
		Type: "appointment_reminder",
		Payload: fiber.Map{
			"appointment_id": appointmentID,
			"patient_name":   patientName,
			"minutes_until":  minutesUntil,
			"timestamp":      time.Now(),
		},
	}
	wsManager.broadcast(message)
}

// Send appointment update notification
func SendAppointmentUpdate(appointmentID uint, patientName string, updateType string) {
	message := WSMessage{
		Type: "appointment_update",
		Payload: fiber.Map{
			"appointment_id": appointmentID,
			"patient_name":   patientName,
			"update_type":    updateType,
			"timestamp":      time.Now(),
		},
	}
	wsManager.broadcast(message)
}

// Send patient update notification
func SendPatientUpdate(patientID uint, patientName string, updateType string) {
	message := WSMessage{
		Type: "patient_update",
		Payload: fiber.Map{
			"patient_id":   patientID,
			"patient_name": patientName,
			"update_type":  updateType,
			"timestamp":    time.Now(),
		},
	}
	wsManager.broadcast(message)
}

// Send general notification
func SendNotification(title, message, notificationType string, userID uint) {
	notification := WSMessage{
		Type: "notification",
		Payload: fiber.Map{
			"title":     title,
			"message":   message,
			"type":      notificationType,
			"user_id":   userID,
			"timestamp": time.Now(),
		},
	}
	wsManager.broadcastToUser(userID, notification)
}

// Send system alert
func SendSystemAlert(title, message, alertType string) {
	alert := WSMessage{
		Type: "system_alert",
		Payload: fiber.Map{
			"title":     title,
			"message":   message,
			"type":      alertType,
			"timestamp": time.Now(),
		},
	}
	wsManager.broadcast(alert)
}
