package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"dentika/server/database"
	"dentika/server/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// WebSocket connection manager with user tracking
type WSManager struct {
	clients   map[*websocket.Conn]*WSClient
	clientsMu sync.RWMutex
}

type WSClient struct {
	UserID   uint
	User     models.User
	Conn     *websocket.Conn
	JoinedAt time.Time
}

var wsManager = &WSManager{
	clients: make(map[*websocket.Conn]*WSClient),
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
func (wm *WSManager) addClient(conn *websocket.Conn, userID uint, user models.User) {
	wm.clientsMu.Lock()
	defer wm.clientsMu.Unlock()

	client := &WSClient{
		UserID:   userID,
		User:     user,
		Conn:     conn,
		JoinedAt: time.Now(),
	}

	wm.clients[conn] = client
	log.Printf("WebSocket client connected for user %d (%s). Total clients: %d", userID, user.Username, len(wm.clients))
}

// Remove client from manager
func (wm *WSManager) removeClient(conn *websocket.Conn) {
	wm.clientsMu.Lock()
	defer wm.clientsMu.Unlock()

	if client, exists := wm.clients[conn]; exists {
		delete(wm.clients, conn)
		log.Printf("WebSocket client disconnected for user %d (%s). Total clients: %d", client.UserID, client.User.Username, len(wm.clients))
	}
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

	for conn, client := range wm.clients {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Failed to send message to client (user %d): %v", client.UserID, err)
			// Remove broken connection
			go wm.removeClient(conn)
		}
	}
}

// Broadcast to specific user
func (wm *WSManager) broadcastToUser(userID uint, message WSMessage) {
	wm.clientsMu.RLock()
	defer wm.clientsMu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}

	for conn, client := range wm.clients {
		if client.UserID == userID {
			if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Failed to send message to user %d: %v", userID, err)
				// Remove broken connection
				go wm.removeClient(conn)
			}
		}
	}
}

// Broadcast to users in specific clinic
func (wm *WSManager) broadcastToClinic(clinicID uint, message WSMessage) {
	wm.clientsMu.RLock()
	defer wm.clientsMu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}

	for conn, client := range wm.clients {
		if client.User.ClinicID != nil && *client.User.ClinicID == clinicID {
			if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Failed to send message to clinic %d user %d: %v", clinicID, client.UserID, err)
				// Remove broken connection
				go wm.removeClient(conn)
			}
		}
	}
}

// Get connected users count
func (wm *WSManager) getConnectedUsers() map[uint]int {
	wm.clientsMu.RLock()
	defer wm.clientsMu.RUnlock()

	userCount := make(map[uint]int)
	for _, client := range wm.clients {
		userCount[client.UserID]++
	}
	return userCount
}

// WebSocket handler
func WebSocketHandler(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		var authenticatedUser *models.User
		var userID uint

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
					user, id, err := handleAuth(conn, wsMsg)
					if err == nil && user != nil {
						authenticatedUser = user
						userID = id
						// Add authenticated client to manager
						wsManager.addClient(conn, userID, *authenticatedUser)
						defer wsManager.removeClient(conn)
					}
				case "ping":
					handlePing(conn)
				default:
					if authenticatedUser == nil {
						conn.WriteJSON(fiber.Map{
							"type":    "error",
							"message": "Not authenticated",
						})
					} else {
						log.Printf("Unknown message type: %s", wsMsg.Type)
					}
				}
			}
		}
	})(c)
}

// Handle authentication - returns user and userID if successful
func handleAuth(conn *websocket.Conn, msg WSMessage) (*models.User, uint, error) {
	if msg.Token == "" {
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "No token provided",
		})
		return nil, 0, fiber.NewError(401, "No token provided")
	}

	// Validate token against database
	var authToken models.AuthToken
	if err := database.DB.Preload("User").Where("token = ?", msg.Token).First(&authToken).Error; err != nil {
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "Invalid token",
		})
		return nil, 0, fiber.NewError(401, "Invalid token")
	}

	// Check if token is expired
	if authToken.IsExpired() {
		database.DB.Delete(&authToken)
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "Token expired",
		})
		return nil, 0, fiber.NewError(401, "Token expired")
	}

	// Check if user is active
	if !authToken.User.IsActive {
		conn.WriteJSON(fiber.Map{
			"type":    "auth_error",
			"message": "Account is inactive",
		})
		return nil, 0, fiber.NewError(403, "Account is inactive")
	}

	// Send success response
	conn.WriteJSON(fiber.Map{
		"type":    "auth_success",
		"message": "Authenticated successfully",
		"user": fiber.Map{
			"id":       authToken.User.ID,
			"username": authToken.User.Username,
			"name":     authToken.User.GetDisplayName(),
		},
	})

	return &authToken.User, authToken.UserID, nil
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

// Send general notification to specific user
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

// Send notification to all users in a clinic
func SendClinicNotification(title, message, notificationType string, clinicID uint) {
	notification := WSMessage{
		Type: "clinic_notification",
		Payload: fiber.Map{
			"title":     title,
			"message":   message,
			"type":      notificationType,
			"clinic_id": clinicID,
			"timestamp": time.Now(),
		},
	}
	wsManager.broadcastToClinic(clinicID, notification)
}

// Send system-wide notification to all connected users
func SendSystemNotification(title, message, notificationType string) {
	notification := WSMessage{
		Type: "system_notification",
		Payload: fiber.Map{
			"title":     title,
			"message":   message,
			"type":      notificationType,
			"timestamp": time.Now(),
		},
	}
	wsManager.broadcast(notification)
}

// Get WebSocket connection statistics
func GetWSStats() fiber.Map {
	userCount := wsManager.getConnectedUsers()
	return fiber.Map{
		"total_connections": len(wsManager.clients),
		"unique_users":      len(userCount),
		"user_connections":  userCount,
	}
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
