package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Notification represents a system notification
type Notification struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Type      NotificationType `json:"type" gorm:"type:varchar(50);not null"`
	Icon      string    `json:"icon" gorm:"size:50"`
	Color     string    `json:"color" gorm:"size:20"`

	// Targeting
	UserID   *uint `json:"user_id,omitempty" gorm:"index"` // Specific user (null for clinic/system wide)
	ClinicID *uint `json:"clinic_id,omitempty" gorm:"index"` // Specific clinic (null for system wide)

	// Relationships
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Clinic *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`

	// Content and actions
	Data    *string `json:"data,omitempty" gorm:"type:json"` // Additional data as JSON
	Actions *string `json:"actions,omitempty" gorm:"type:json"` // Actions as JSON

	// Scheduling
	ScheduledFor *time.Time `json:"scheduled_for,omitempty" gorm:"index"` // For delayed notifications
	ExpiresAt    *time.Time `json:"expires_at,omitempty" gorm:"index"`    // When notification should expire

	// Metadata
	CreatedByID *uint `json:"created_by_id,omitempty" gorm:"index"`
	CreatedBy   *User `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationTypeSuccess            NotificationType = "success"
	NotificationTypeError              NotificationType = "error"
	NotificationTypeWarning            NotificationType = "warning"
	NotificationTypeInfo               NotificationType = "info"
	NotificationTypeAppointmentReminder NotificationType = "appointment_reminder"
	NotificationTypeAppointmentUpdate   NotificationType = "appointment_update"
	NotificationTypePatientUpdate       NotificationType = "patient_update"
	NotificationTypeSystemAlert         NotificationType = "system_alert"
	NotificationTypeClinicAnnouncement  NotificationType = "clinic_announcement"
	NotificationTypePeerReviewUpdate    NotificationType = "peer_review_update"
	NotificationTypeInventoryAlert      NotificationType = "inventory_alert"
)

// NotificationRecipient tracks which users have read/removed a notification
type NotificationRecipient struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	NotificationID uint         `json:"notification_id" gorm:"not null;index"`
	Notification   Notification `json:"notification" gorm:"foreignKey:NotificationID"`

	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// Status - simplified to only track read state
	IsRead bool       `json:"is_read" gorm:"default:false;index"`
	ReadAt *time.Time `json:"read_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // For soft delete (removal)
}

// NotificationAction represents an action that can be taken on a notification
type NotificationAction struct {
	Label    string            `json:"label"`
	Action   string            `json:"action"`
	URL      string            `json:"url,omitempty"`
	Method   string            `json:"method,omitempty"` // GET, POST, etc.
	Payload  map[string]interface{} `json:"payload,omitempty"`
	Style    string            `json:"style,omitempty"` // primary, secondary, danger
}

// Scope defines the scope of a notification
type NotificationScope string

const (
	ScopeUser   NotificationScope = "user"   // Single user
	ScopeClinic NotificationScope = "clinic" // All users in a clinic
	ScopeSystem NotificationScope = "system" // All users system-wide
)

// GetScope returns the scope of the notification
func (n *Notification) GetScope() NotificationScope {
	if n.UserID != nil {
		return ScopeUser
	}
	if n.ClinicID != nil {
		return ScopeClinic
	}
	return ScopeSystem
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	return n.ExpiresAt != nil && n.ExpiresAt.Before(time.Now())
}

// IsScheduled checks if the notification is scheduled for future delivery
func (n *Notification) IsScheduled() bool {
	return n.ScheduledFor != nil && n.ScheduledFor.After(time.Now())
}

// GetNATSSubject returns the appropriate NATS subject for this notification
func (n *Notification) GetNATSSubject() string {
	switch n.GetScope() {
	case ScopeUser:
		return "dentika.user." + string(rune(*n.UserID)) + ".notifications"
	case ScopeClinic:
		return "dentika.clinic." + string(rune(*n.ClinicID)) + ".notifications"
	case ScopeSystem:
		return "dentika.system.notifications"
	default:
		return "dentika.system.notifications"
	}
}

// GetDefaultIcon returns a default icon based on notification type
func (n *Notification) GetDefaultIcon() string {
	if n.Icon != "" {
		return n.Icon
	}

	iconMap := map[NotificationType]string{
		NotificationTypeSuccess:            "check-circle",
		NotificationTypeError:              "x-circle",
		NotificationTypeWarning:            "exclamation-triangle",
		NotificationTypeInfo:               "information-circle",
		NotificationTypeAppointmentReminder: "clock",
		NotificationTypeAppointmentUpdate:   "calendar",
		NotificationTypePatientUpdate:       "user",
		NotificationTypeSystemAlert:         "exclamation-triangle",
		NotificationTypeClinicAnnouncement:  "megaphone",
		NotificationTypePeerReviewUpdate:    "users",
		NotificationTypeInventoryAlert:      "box",
	}

	if icon, exists := iconMap[n.Type]; exists {
		return icon
	}
	return "information-circle"
}

// GetDefaultColor returns a default color based on notification type
func (n *Notification) GetDefaultColor() string {
	if n.Color != "" {
		return n.Color
	}

	colorMap := map[NotificationType]string{
		NotificationTypeSuccess:            "green",
		NotificationTypeError:              "red",
		NotificationTypeWarning:            "yellow",
		NotificationTypeInfo:               "blue",
		NotificationTypeAppointmentReminder: "blue",
		NotificationTypeAppointmentUpdate:   "purple",
		NotificationTypePatientUpdate:       "green",
		NotificationTypeSystemAlert:         "red",
		NotificationTypeClinicAnnouncement:  "blue",
		NotificationTypePeerReviewUpdate:    "purple",
		NotificationTypeInventoryAlert:      "yellow",
	}

	if color, exists := colorMap[n.Type]; exists {
		return color
	}
	return "blue"
}

// Custom JSON marshaling to handle Actions and Data fields properly
type NotificationJSON struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      NotificationType `json:"type"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	UserID    *uint     `json:"user_id,omitempty"`
	ClinicID  *uint     `json:"clinic_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Actions   []NotificationAction `json:"actions,omitempty"`
	ScheduledFor *time.Time `json:"scheduled_for,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedByID  *uint     `json:"created_by_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MarshalJSON customizes JSON output to parse Actions and Data strings
func (n *Notification) MarshalJSON() ([]byte, error) {
	notificationJSON := NotificationJSON{
		ID:          n.ID,
		Title:       n.Title,
		Message:     n.Message,
		Type:        n.Type,
		Icon:        n.Icon,
		Color:       n.Color,
		UserID:      n.UserID,
		ClinicID:    n.ClinicID,
		ScheduledFor: n.ScheduledFor,
		ExpiresAt:   n.ExpiresAt,
		CreatedByID: n.CreatedByID,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}

	// Parse Data field if it exists
	if n.Data != nil && *n.Data != "" {
		var dataObj interface{}
		if err := json.Unmarshal([]byte(*n.Data), &dataObj); err == nil {
			notificationJSON.Data = dataObj
		} else {
			notificationJSON.Data = *n.Data
		}
	}

	// Parse Actions field if it exists
	if n.Actions != nil && *n.Actions != "" {
		var actionsArray []NotificationAction
		if err := json.Unmarshal([]byte(*n.Actions), &actionsArray); err == nil {
			notificationJSON.Actions = actionsArray
		}
	}

	return json.Marshal(notificationJSON)
}