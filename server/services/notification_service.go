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

// createRecipientsAndDistribute creates recipient records and distributes via NATS
func (ns *NotificationService) createRecipientsAndDistribute(notification *models.Notification) error {
	var recipients []models.User

	// Determine recipients based on scope
	switch notification.GetScope() {
	case models.ScopeUser:
		// Single user
		var user models.User
		if err := ns.db.First(&user, *notification.UserID).Error; err != nil {
			return fmt.Errorf("failed to find user: %v", err)
		}
		recipients = []models.User{user}

	case models.ScopeClinic:
		// All users in clinic
		if err := ns.db.Where("clinic_id = ?", *notification.ClinicID).Find(&recipients).Error; err != nil {
			return fmt.Errorf("failed to find clinic users: %v", err)
		}

	case models.ScopeSystem:
		// All active users
		if err := ns.db.Find(&recipients).Error; err != nil {
			return fmt.Errorf("failed to find all users: %v", err)
		}
	}

	// Create recipient records
	var recipientRecords []models.NotificationRecipient
	for _, user := range recipients {
		recipientRecords = append(recipientRecords, models.NotificationRecipient{
			NotificationID: notification.ID,
			UserID:         user.ID,
			DeliveredAt:    time.Now(),
		})
	}

	if len(recipientRecords) > 0 {
		if err := ns.db.Create(&recipientRecords).Error; err != nil {
			return fmt.Errorf("failed to create recipient records: %v", err)
		}
	}

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
	now := time.Now()

	return ns.db.Model(&models.NotificationRecipient{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

// MarkAsDismissed marks a notification as dismissed for a specific user
func (ns *NotificationService) MarkAsDismissed(notificationID, userID uint) error {
	now := time.Now()

	return ns.db.Model(&models.NotificationRecipient{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_dismissed": true,
			"dismissed_at": &now,
		}).Error
}

// GetUserNotifications retrieves notifications for a specific user
func (ns *NotificationService) GetUserNotifications(userID uint, limit, offset int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	// Base query - get notifications that are either:
	// 1. User-specific (user_id = userID)
	// 2. Clinic-wide (clinic_id = user's clinic_id)
	// 3. System-wide (user_id IS NULL AND clinic_id IS NULL)

	// First get the user's clinic_id
	var user models.User
	if err := ns.db.First(&user, userID).Error; err != nil {
		return nil, 0, err
	}

	query := ns.db.Model(&models.Notification{}).
		Joins("LEFT JOIN notification_recipients nr ON notifications.id = nr.notification_id AND nr.user_id = ?", userID).
		Where("(notifications.user_id = ? OR (notifications.clinic_id = ? AND notifications.user_id IS NULL) OR (notifications.user_id IS NULL AND notifications.clinic_id IS NULL)) AND (notifications.scheduled_for IS NULL OR notifications.scheduled_for <= ?) AND (notifications.expires_at IS NULL OR notifications.expires_at > ?) AND (nr.is_dismissed = false OR nr.is_dismissed IS NULL)",
			userID, user.ClinicID, time.Now(), time.Now()).
		Order("notifications.created_at DESC")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with recipient info
	err := query.
		Select("notifications.*, nr.is_read, nr.read_at, nr.is_dismissed, nr.dismissed_at").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

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

	err := ns.db.Model(&models.Notification{}).
		Joins("LEFT JOIN notification_recipients nr ON notifications.id = nr.notification_id AND nr.user_id = ?", userID).
		Where("(notifications.user_id = ? OR (notifications.clinic_id = ? AND notifications.user_id IS NULL) OR (notifications.user_id IS NULL AND notifications.clinic_id IS NULL)) AND (nr.is_read = false OR nr.is_read IS NULL) AND (nr.is_dismissed = false OR nr.is_dismissed IS NULL) AND (notifications.scheduled_for IS NULL OR notifications.scheduled_for <= ?) AND (notifications.expires_at IS NULL OR notifications.expires_at > ?)",
			userID, user.ClinicID, time.Now(), time.Now()).
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