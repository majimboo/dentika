package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type NotificationService struct {
	db   *gorm.DB
	nats *nats.Conn
}

func NewNotificationService(natsConn *nats.Conn) *NotificationService {
	return &NotificationService{
		db:   database.DB,
		nats: natsConn,
	}
}

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	Title       string                     `json:"title"`
	Message     string                     `json:"message"`
	Type        models.NotificationType    `json:"type"`
	Icon        string                     `json:"icon,omitempty"`
	Color       string                     `json:"color,omitempty"`
	UserID      *uint                      `json:"user_id,omitempty"`
	ClinicID    *uint                      `json:"clinic_id,omitempty"`
	Data        map[string]interface{}     `json:"data,omitempty"`
	Actions     []models.NotificationAction `json:"actions,omitempty"`
	ScheduledFor *time.Time                `json:"scheduled_for,omitempty"`
	ExpiresAt   *time.Time                 `json:"expires_at,omitempty"`
	CreatedByID *uint                      `json:"created_by_id,omitempty"`
}

// CreateNotification creates a new notification and distributes it
func (ns *NotificationService) CreateNotification(req CreateNotificationRequest) (*models.Notification, error) {
	log.Printf("CreateNotification called with title: %s, type: %s", req.Title, req.Type)
	// Prepare data JSON
	var dataJSON *string
	if req.Data != nil {
		dataBytes, err := json.Marshal(req.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %v", err)
		}
		dataStr := string(dataBytes)
		// Only set if not empty JSON
		if dataStr != "null" && dataStr != "{}" && dataStr != "[]" {
			dataJSON = &dataStr
		}
	}

	// Prepare actions JSON
	var actionsJSON *string
	if req.Actions != nil {
		actionsBytes, err := json.Marshal(req.Actions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal actions: %v", err)
		}
		actionsStr := string(actionsBytes)
		// Only set if not empty JSON
		if actionsStr != "null" && actionsStr != "{}" && actionsStr != "[]" {
			actionsJSON = &actionsStr
		}
	}

	// Create notification
	notification := models.Notification{
		Title:        req.Title,
		Message:      req.Message,
		Type:         req.Type,
		Icon:         req.Icon,
		Color:        req.Color,
		UserID:       req.UserID,
		ClinicID:     req.ClinicID,
		Data:         dataJSON,
		Actions:      actionsJSON,
		ScheduledFor: req.ScheduledFor,
		ExpiresAt:    req.ExpiresAt,
		CreatedByID:  req.CreatedByID,
	}

	// Set defaults
	if notification.Icon == "" {
		notification.Icon = notification.GetDefaultIcon()
	}
	if notification.Color == "" {
		notification.Color = notification.GetDefaultColor()
	}

	// Save to database
	log.Printf("Attempting to save notification to database...")
	if err := ns.db.Create(&notification).Error; err != nil {
		log.Printf("ERROR saving notification to database: %v", err)
		return nil, fmt.Errorf("failed to create notification: %v", err)
	}
	log.Printf("Successfully saved notification to database with ID: %d", notification.ID)

	// Create recipient records and distribute
	if err := ns.createRecipientsAndDistribute(&notification); err != nil {
		log.Printf("Failed to distribute notification %d: %v", notification.ID, err)
	}

	return &notification, nil
}

// createRecipientsAndDistribute distributes notifications via NATS without creating recipient records
func (ns *NotificationService) createRecipientsAndDistribute(notification *models.Notification) error {
	// Skip creating recipient records - they will be created on-demand when users fetch notifications

	// Distribute via NATS if not scheduled
	if notification.ScheduledFor == nil || notification.ScheduledFor.Before(time.Now()) {
		return ns.distributeNotification(notification)
	}

	return nil
}

// distributeNotification sends the notification via NATS
func (ns *NotificationService) distributeNotification(notification *models.Notification) error {
	if ns.nats == nil {
		return fmt.Errorf("NATS connection not available")
	}

	// Prepare notification data for NATS
	natsData := map[string]interface{}{
		"id":         notification.ID,
		"title":      notification.Title,
		"message":    notification.Message,
		"type":       notification.Type,
		"icon":       notification.GetDefaultIcon(),
		"color":      notification.GetDefaultColor(),
		"created_at": notification.CreatedAt,
	}

	// Add optional data and actions if present
	if notification.Data != nil {
		var dataObj interface{}
		if err := json.Unmarshal([]byte(*notification.Data), &dataObj); err == nil {
			natsData["data"] = dataObj
		} else {
			natsData["data"] = *notification.Data
		}
	}
	if notification.Actions != nil {
		var actionsArray []models.NotificationAction
		if err := json.Unmarshal([]byte(*notification.Actions), &actionsArray); err == nil {
			natsData["actions"] = actionsArray
		} else {
			natsData["actions"] = *notification.Actions
		}
	}

	// Serialize data
	data, err := json.Marshal(natsData)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %v", err)
	}

	// Determine NATS subjects and publish
	subjects := ns.getNATSSubjects(notification)
	for _, subject := range subjects {
		if err := ns.nats.Publish(subject, data); err != nil {
			log.Printf("Failed to publish notification to %s: %v", subject, err)
		} else {
			log.Printf("Published notification %d to %s", notification.ID, subject)
		}
	}

	return nil
}

// getNATSSubjects returns all NATS subjects for the notification
func (ns *NotificationService) getNATSSubjects(notification *models.Notification) []string {
	var subjects []string

	switch notification.GetScope() {
	case models.ScopeUser:
		subjects = append(subjects, fmt.Sprintf("dentika.user.%d.notifications", *notification.UserID))

	case models.ScopeClinic:
		subjects = append(subjects, fmt.Sprintf("dentika.clinic.%d.notifications", *notification.ClinicID))

	case models.ScopeSystem:
		subjects = append(subjects, "dentika.system.notifications")
	}

	return subjects
}

// MarkAsRead marks a notification as read for a specific user
func (ns *NotificationService) MarkAsRead(notificationID, userID uint) error {
	log.Printf("SERVICE: MarkAsRead called for notification %d, user %d", notificationID, userID)

	// Check if recipient record already exists
	var existingRecipient models.NotificationRecipient
	err := ns.db.Where("notification_id = ? AND user_id = ?", notificationID, userID).First(&existingRecipient).Error

	if err == gorm.ErrRecordNotFound {
		log.Printf("SERVICE: No existing recipient record found, creating new one")
		// Create new read record
		now := time.Now()
		recipient := models.NotificationRecipient{
			NotificationID: notificationID,
			UserID:         userID,
			IsRead:         true,
			ReadAt:         &now,
		}
		err := ns.db.Create(&recipient).Error
		if err != nil {
			log.Printf("SERVICE: ERROR creating recipient record: %v", err)
			return err
		}
		log.Printf("SERVICE: Successfully created recipient record with ID %d", recipient.ID)
		return nil
	}

	if err != nil {
		log.Printf("SERVICE: ERROR querying existing recipient: %v", err)
		return err
	}

	log.Printf("SERVICE: Found existing recipient record ID %d, updating it", existingRecipient.ID)
	// Update existing record to mark as read
	now := time.Now()
	err = ns.db.Model(&existingRecipient).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": &now,
	}).Error

	if err != nil {
		log.Printf("SERVICE: ERROR updating recipient record: %v", err)
		return err
	}

	log.Printf("SERVICE: Successfully updated recipient record")
	return nil
}

// MarkAsDismissed permanently removes a notification for a specific user (soft delete)
func (ns *NotificationService) MarkAsDismissed(notificationID, userID uint) error {
	log.Printf("SERVICE: MarkAsDismissed called for notification %d, user %d", notificationID, userID)

	// Find the recipient record
	var recipient models.NotificationRecipient
	err := ns.db.Where("notification_id = ? AND user_id = ?", notificationID, userID).First(&recipient).Error

	if err == gorm.ErrRecordNotFound {
		log.Printf("SERVICE: No recipient record exists, creating one to soft delete")
		// No recipient record exists - create one and immediately soft delete it
		now := time.Now()
		recipient = models.NotificationRecipient{
			NotificationID: notificationID,
			UserID:         userID,
			IsRead:         true,  // Mark as read when removing
			ReadAt:         &now,
		}
		// Create and immediately soft delete
		if err := ns.db.Create(&recipient).Error; err != nil {
			log.Printf("SERVICE: ERROR creating recipient record: %v", err)
			return err
		}
		log.Printf("SERVICE: Created recipient record with ID %d", recipient.ID)
	} else if err != nil {
		log.Printf("SERVICE: ERROR querying recipient record: %v", err)
		return err
	}

	// Soft delete the recipient record
	err = ns.db.Delete(&recipient).Error
	if err != nil {
		log.Printf("SERVICE: ERROR soft deleting recipient record: %v", err)
		return err
	}

	log.Printf("SERVICE: Successfully soft deleted recipient record")
	return nil
}


// NotificationWithStatus represents a notification with read status
type NotificationWithStatus struct {
	ID           uint                         `json:"id" gorm:"primarykey"`
	Title        string                       `json:"title"`
	Message      string                       `json:"message"`
	Type         models.NotificationType      `json:"type"`
	Icon         string                       `json:"icon"`
	Color        string                       `json:"color"`
	UserID       *uint                        `json:"user_id,omitempty"`
	ClinicID     *uint                        `json:"clinic_id,omitempty"`
	Data         *string                      `json:"data,omitempty"`
	Actions      *string                      `json:"actions,omitempty"`
	ScheduledFor *time.Time                   `json:"scheduled_for,omitempty"`
	ExpiresAt    *time.Time                   `json:"expires_at,omitempty"`
	CreatedByID  *uint                        `json:"created_by_id,omitempty"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
	IsRead       bool                         `json:"is_read"`
	ReadAt       *time.Time                   `json:"read_at,omitempty"`
}

// MarshalJSON customizes JSON output for NotificationWithStatus
func (nws *NotificationWithStatus) MarshalJSON() ([]byte, error) {
	type Alias NotificationWithStatus

	// Create an alias to avoid infinite recursion
	aux := &struct {
		*Alias
		Data    interface{}                  `json:"data,omitempty"`
		Actions []models.NotificationAction  `json:"actions,omitempty"`
	}{
		Alias: (*Alias)(nws),
	}

	// Parse Data field if it exists
	if nws.Data != nil && *nws.Data != "" {
		var dataObj interface{}
		if err := json.Unmarshal([]byte(*nws.Data), &dataObj); err == nil {
			aux.Data = dataObj
		} else {
			aux.Data = *nws.Data
		}
	}

	// Parse Actions field if it exists
	if nws.Actions != nil && *nws.Actions != "" {
		var actionsArray []models.NotificationAction
		if err := json.Unmarshal([]byte(*nws.Actions), &actionsArray); err == nil {
			aux.Actions = actionsArray
		}
	}

	return json.Marshal(aux)
}

// GetUserNotifications retrieves notifications for a specific user
func (ns *NotificationService) GetUserNotifications(userID uint, limit, offset int) ([]NotificationWithStatus, int64, error) {
	var notifications []NotificationWithStatus
	var total int64

	// First get the user's clinic_id
	var user models.User
	if err := ns.db.First(&user, userID).Error; err != nil {
		return nil, 0, err
	}

	// Base query - get all relevant notifications:
	// 1. User-specific (user_id = userID)
	// 2. Clinic-wide (clinic_id = user's clinic_id)
	// 3. System-wide (user_id IS NULL AND clinic_id IS NULL)
	// EXCLUDE only notifications that have been REMOVED (soft-deleted)
	// INCLUDE notifications that are read but not removed
	baseQuery := ns.db.Model(&models.Notification{}).
		Joins("LEFT JOIN notification_recipients nr ON notifications.id = nr.notification_id AND nr.user_id = ? AND nr.deleted_at IS NULL", userID).
		Where("(notifications.user_id = ? OR (notifications.clinic_id = ? AND notifications.user_id IS NULL) OR (notifications.user_id IS NULL AND notifications.clinic_id IS NULL)) AND (notifications.scheduled_for IS NULL OR notifications.scheduled_for <= ?) AND (notifications.expires_at IS NULL OR notifications.expires_at > ?)",
			userID, user.ClinicID, time.Now(), time.Now()).
		Where("NOT EXISTS (SELECT 1 FROM notification_recipients nr_removed WHERE nr_removed.notification_id = notifications.id AND nr_removed.user_id = ? AND nr_removed.deleted_at IS NOT NULL)", userID)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated notifications with recipient status
	err := baseQuery.
		Select("notifications.*, COALESCE(nr.is_read, false) as is_read, nr.read_at").
		Order("notifications.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, err
}

// GetUnreadCount returns the count of unread notifications for a user
func (ns *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64

	// Get user's clinic_id
	var user models.User
	if err := ns.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	// Count notifications that are:
	// 1. Relevant to the user (user/clinic/system scope)
	// 2. Not removed (no soft-deleted recipient record)
	// 3. Not read (no recipient record with is_read=true OR recipient record doesn't exist)
	err := ns.db.Model(&models.Notification{}).
		Joins("LEFT JOIN notification_recipients nr ON notifications.id = nr.notification_id AND nr.user_id = ? AND nr.deleted_at IS NULL", userID).
		Where("(notifications.user_id = ? OR (notifications.clinic_id = ? AND notifications.user_id IS NULL) OR (notifications.user_id IS NULL AND notifications.clinic_id IS NULL)) AND (notifications.scheduled_for IS NULL OR notifications.scheduled_for <= ?) AND (notifications.expires_at IS NULL OR notifications.expires_at > ?)",
			userID, user.ClinicID, time.Now(), time.Now()).
		Where("NOT EXISTS (SELECT 1 FROM notification_recipients nr_removed WHERE nr_removed.notification_id = notifications.id AND nr_removed.user_id = ? AND nr_removed.deleted_at IS NOT NULL)", userID).
		Where("(nr.id IS NULL OR nr.is_read = false)").
		Count(&count).Error

	return count, err
}

// ProcessScheduledNotifications processes notifications that are scheduled for delivery
func (ns *NotificationService) ProcessScheduledNotifications() error {
	var notifications []models.Notification

	err := ns.db.Where("scheduled_for <= ? AND scheduled_for IS NOT NULL", time.Now()).Find(&notifications).Error
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := ns.distributeNotification(&notification); err != nil {
			log.Printf("Failed to distribute scheduled notification %d: %v", notification.ID, err)
		} else {
			// Clear scheduled_for to prevent re-processing
			ns.db.Model(&notification).Update("scheduled_for", nil)
		}
	}

	return nil
}

// CleanupExpiredNotifications removes expired notifications
func (ns *NotificationService) CleanupExpiredNotifications() error {
	return ns.db.Where("expires_at <= ?", time.Now()).Delete(&models.Notification{}).Error
}

// Convenience methods for common notification types

func (ns *NotificationService) CreateAppointmentReminder(appointment models.Appointment, minutesUntil int) error {
	_, err := ns.CreateNotification(CreateNotificationRequest{
		Title:   "Appointment Reminder",
		Message: fmt.Sprintf("%s %s has an appointment in %d minutes", appointment.Patient.FirstName, appointment.Patient.LastName, minutesUntil),
		Type:    models.NotificationTypeAppointmentReminder,
		UserID:  &appointment.DoctorID,
		Data: map[string]interface{}{
			"appointment_id": appointment.ID,
			"patient_id":     appointment.PatientID,
			"minutes_until":  minutesUntil,
		},
		Actions: []models.NotificationAction{
			{
				Label:  "View Appointment",
				Action: "view-appointment",
				URL:    fmt.Sprintf("/appointments/%d", appointment.ID),
			},
		},
	})
	return err
}

func (ns *NotificationService) CreateAppointmentUpdate(appointment models.Appointment, updateType string) error {
	var title, message string

	switch updateType {
	case "scheduled":
		title = "New Appointment Scheduled"
		message = fmt.Sprintf("Appointment scheduled for %s %s", appointment.Patient.FirstName, appointment.Patient.LastName)
	case "cancelled":
		title = "Appointment Cancelled"
		message = fmt.Sprintf("Appointment with %s %s has been cancelled", appointment.Patient.FirstName, appointment.Patient.LastName)
	case "rescheduled":
		title = "Appointment Rescheduled"
		message = fmt.Sprintf("Appointment with %s %s has been rescheduled", appointment.Patient.FirstName, appointment.Patient.LastName)
	default:
		title = "Appointment Updated"
		message = fmt.Sprintf("Appointment with %s %s has been updated", appointment.Patient.FirstName, appointment.Patient.LastName)
	}

	_, err := ns.CreateNotification(CreateNotificationRequest{
		Title:    title,
		Message:  message,
		Type:     models.NotificationTypeAppointmentUpdate,
		ClinicID: &appointment.ClinicID,
		Data: map[string]interface{}{
			"appointment_id": appointment.ID,
			"patient_id":     appointment.PatientID,
			"update_type":    updateType,
		},
		Actions: []models.NotificationAction{
			{
				Label:  "View Appointment",
				Action: "view-appointment",
				URL:    fmt.Sprintf("/appointments/%d", appointment.ID),
			},
		},
	})
	return err
}

func (ns *NotificationService) CreatePatientUpdate(patient models.Patient, updateType string) error {
	var title, message string

	switch updateType {
	case "created":
		title = "New Patient Added"
		message = fmt.Sprintf("New patient %s %s has been added", patient.FirstName, patient.LastName)
	default:
		title = "Patient Updated"
		message = fmt.Sprintf("Patient %s %s has been updated", patient.FirstName, patient.LastName)
	}

	_, err := ns.CreateNotification(CreateNotificationRequest{
		Title:    title,
		Message:  message,
		Type:     models.NotificationTypePatientUpdate,
		ClinicID: &patient.ClinicID,
		Data: map[string]interface{}{
			"patient_id":  patient.ID,
			"update_type": updateType,
		},
		Actions: []models.NotificationAction{
			{
				Label:  "View Patient",
				Action: "view-patient",
				URL:    fmt.Sprintf("/patients/%d", patient.ID),
			},
		},
	})
	return err
}

func (ns *NotificationService) CreateSystemAlert(title, message string, severity models.NotificationType) error {
	_, err := ns.CreateNotification(CreateNotificationRequest{
		Title:   title,
		Message: message,
		Type:    severity,
	})
	return err
}