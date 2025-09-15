package handlers

import (
	"dentika/server/models"
	"dentika/server/services"
	"fmt"
	"log"
)

// Replacement functions for websocket notifications using the new notification service

// SendAppointmentUpdate creates an appointment update notification
func SendAppointmentUpdate(appointmentID uint, patientName string, updateType string, clinicID uint) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		var title, message string
		var notificationType models.NotificationType

		switch updateType {
		case "scheduled":
			title = "New Appointment Scheduled"
			message = fmt.Sprintf("Appointment scheduled for %s", patientName)
			notificationType = models.NotificationTypeAppointmentUpdate
		case "cancelled":
			title = "Appointment Cancelled"
			message = fmt.Sprintf("Appointment with %s has been cancelled", patientName)
			notificationType = models.NotificationTypeAppointmentUpdate
		case "updated":
			title = "Appointment Updated"
			message = fmt.Sprintf("Appointment with %s has been updated", patientName)
			notificationType = models.NotificationTypeAppointmentUpdate
		default:
			title = "Appointment Status Changed"
			message = fmt.Sprintf("Appointment with %s status: %s", patientName, updateType)
			notificationType = models.NotificationTypeAppointmentUpdate
		}

		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:    title,
			Message:  message,
			Type:     notificationType,
			ClinicID: &clinicID,
			Data: map[string]interface{}{
				"appointment_id": appointmentID,
				"patient_name":   patientName,
				"update_type":    updateType,
			},
			Actions: []models.NotificationAction{
				{
					Label:  "View Appointment",
					Action: "view-appointment",
					URL:    fmt.Sprintf("/appointments/%d", appointmentID),
				},
			},
		})

		if err != nil {
			log.Printf("Failed to create appointment update notification: %v", err)
		}
	}()
}

// SendPatientUpdate creates a patient update notification
func SendPatientUpdate(patientID uint, patientName string, updateType string, clinicID uint) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		var title, message string

		switch updateType {
		case "created":
			title = "New Patient Added"
			message = fmt.Sprintf("New patient %s has been added", patientName)
		default:
			title = "Patient Updated"
			message = fmt.Sprintf("Patient %s has been updated", patientName)
		}

		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:    title,
			Message:  message,
			Type:     models.NotificationTypePatientUpdate,
			ClinicID: &clinicID,
			Data: map[string]interface{}{
				"patient_id":   patientID,
				"patient_name": patientName,
				"update_type":  updateType,
			},
			Actions: []models.NotificationAction{
				{
					Label:  "View Patient",
					Action: "view-patient",
					URL:    fmt.Sprintf("/patients/%d", patientID),
				},
			},
		})

		if err != nil {
			log.Printf("Failed to create patient update notification: %v", err)
		}
	}()
}

// SendNotification creates a user-specific notification
func SendNotification(title, message, notificationType string, userID uint) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		var nType models.NotificationType
		switch notificationType {
		case "welcome":
			nType = models.NotificationTypeInfo
		case "task_assigned":
			nType = models.NotificationTypeInfo
		default:
			nType = models.NotificationTypeInfo
		}

		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:   title,
			Message: message,
			Type:    nType,
			UserID:  &userID,
		})

		if err != nil {
			log.Printf("Failed to create user notification: %v", err)
		}
	}()
}

// SendClinicNotification creates a clinic-wide notification
func SendClinicNotification(title, message, notificationType string, clinicID uint) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		var nType models.NotificationType
		switch notificationType {
		case "inventory", "inventory_low", "restock_alert":
			nType = models.NotificationTypeInventoryAlert
		case "appointment_cancelled":
			nType = models.NotificationTypeAppointmentUpdate
		case "patient_created", "patient_registered":
			nType = models.NotificationTypePatientUpdate
		case "system_maintenance":
			nType = models.NotificationTypeInfo
		default:
			nType = models.NotificationTypeInfo
		}

		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:    title,
			Message:  message,
			Type:     nType,
			ClinicID: &clinicID,
		})

		if err != nil {
			log.Printf("Failed to create clinic notification: %v", err)
		}
	}()
}

// SendSystemNotification creates a system-wide notification
func SendSystemNotification(title, message, notificationType string) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		var nType models.NotificationType
		switch notificationType {
		case "system_update":
			nType = models.NotificationTypeInfo
		default:
			nType = models.NotificationTypeInfo
		}

		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:   title,
			Message: message,
			Type:    nType,
		})

		if err != nil {
			log.Printf("Failed to create system notification: %v", err)
		}
	}()
}

// SendSystemAlert creates a system alert
func SendSystemAlert(title, message, alertType string) {
	if notificationService == nil {
		log.Printf("Notification service not initialized")
		return
	}

	go func() {
		_, err := notificationService.CreateNotification(services.CreateNotificationRequest{
			Title:   title,
			Message: message,
			Type:    models.NotificationTypeSystemAlert,
		})

		if err != nil {
			log.Printf("Failed to create system alert: %v", err)
		}
	}()
}

// GetWSStats returns empty stats since we're no longer using WebSocket
func GetWSStats() map[string]interface{} {
	return map[string]interface{}{
		"message":            "WebSocket deprecated - using NATS now",
		"total_connections":  0,
		"unique_users":       0,
		"user_connections":   map[string]int{},
	}
}