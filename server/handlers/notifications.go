package handlers

import (
	"log"
	"time"

	"dentika/server/database"
	"dentika/server/models"
)

// StartAppointmentReminderService starts a background service to send appointment reminders
func StartAppointmentReminderService() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				checkAndSendReminders()
			}
		}
	}()
	log.Println("Appointment reminder service started")
}

// checkAndSendReminders checks for upcoming appointments and sends reminders
func checkAndSendReminders() {
	now := time.Now()
	fifteenMinutesFromNow := now.Add(15 * time.Minute)
	thirtyMinutesFromNow := now.Add(30 * time.Minute)

	// Find appointments starting in 15-30 minutes
	var upcomingAppointments []models.Appointment
	err := database.DB.Preload("Patient").Preload("Branch").
		Where("start_time > ? AND start_time <= ? AND status IN (?, ?)",
			fifteenMinutesFromNow, thirtyMinutesFromNow,
			models.StatusScheduled, models.StatusConfirmed).
		Find(&upcomingAppointments).Error

	if err != nil {
		log.Printf("Error fetching upcoming appointments: %v", err)
		return
	}

	for _, appointment := range upcomingAppointments {
		minutesUntil := int(appointment.StartTime.Sub(now).Minutes())

		// Send reminder via notification service
		if notificationService != nil {
			go notificationService.CreateAppointmentReminder(appointment, minutesUntil)
		}

		log.Printf("Sent reminder for appointment %d in clinic %d (%d minutes until start)", appointment.ID, appointment.Branch.ClinicID, minutesUntil)
	}
}

// SendWelcomeNotification sends a welcome notification to a new user
func SendWelcomeNotification(userID uint, userName string) {
	go SendNotification(
		"Welcome to Dentika!",
		"Welcome "+userName+", your account has been created successfully.",
		"welcome",
		userID,
	)
}

// SendLowInventoryAlert sends an alert to all users in the clinic
func SendLowInventoryAlert(itemName string, currentStock int, minStock int, clinicID uint) {
	go SendClinicNotification(
		"Low Inventory Alert",
		itemName+" is running low (Current: "+string(rune(currentStock))+" | Minimum: "+string(rune(minStock))+")",
		"inventory",
		clinicID,
	)
}

// SendAppointmentCancellationAlert sends alert when appointment is cancelled
func SendAppointmentCancellationAlert(appointmentID uint, patientName string, reason string, clinicID uint) {
	message := "Appointment with " + patientName + " has been cancelled"
	if reason != "" {
		message += ". Reason: " + reason
	}

	// Send to clinic users
	go SendClinicNotification(
		"Appointment Cancelled",
		message,
		"appointment_cancelled",
		clinicID,
	)
}

// SendPatientCreatedNotification notifies clinic users about new patient
func SendPatientCreatedNotification(patientID uint, patientName string, clinicID uint) {
	// Send both clinic notification and patient update
	go SendClinicNotification(
		"New Patient Added",
		"Patient "+patientName+" has been added to the system.",
		"patient_created",
		clinicID,
	)
	go SendPatientUpdate(patientID, patientName, "created", clinicID)
}

// SendInventoryRestockAlert notifies when items need restocking
func SendInventoryRestockAlert(itemName string, clinicID uint) {
	go SendClinicNotification(
		"Restock Required",
		itemName+" needs to be restocked.",
		"restock_alert",
		clinicID,
	)
}

// Example usage functions for different notification scenarios

// NotifyClinicAboutNewPatient sends notification when a new patient is added
func NotifyClinicAboutNewPatient(patientName string, clinicID uint) {
	go SendClinicNotification(
		"New Patient Registered",
		"Patient "+patientName+" has been registered in the system.",
		"patient_registered",
		clinicID,
	)
}

// NotifyClinicAboutAppointmentCancellation sends notification when appointment is cancelled
func NotifyClinicAboutAppointmentCancellation(patientName, doctorName string, clinicID uint) {
	go SendClinicNotification(
		"Appointment Cancelled",
		"Appointment with "+patientName+" and Dr. "+doctorName+" has been cancelled.",
		"appointment_cancelled",
		clinicID,
	)
}

// NotifyClinicAboutLowInventory sends inventory alert to clinic
func NotifyClinicAboutLowInventory(itemName string, currentStock int, clinicID uint) {
	go SendClinicNotification(
		"⚠️ Low Inventory Alert",
		itemName+" stock is low ("+string(rune(currentStock))+" remaining). Please restock soon.",
		"inventory_low",
		clinicID,
	)
}

// NotifyUserAboutTaskAssignment sends personal notification to specific user
func NotifyUserAboutTaskAssignment(userID uint, taskTitle string) {
	go SendNotification(
		"New Task Assigned",
		"You have been assigned: "+taskTitle,
		"task_assigned",
		userID,
	)
}

// NotifyClinicAboutSystemMaintenance sends system-wide maintenance notification
func NotifyClinicAboutSystemMaintenance(message string, clinicID uint) {
	go SendClinicNotification(
		"🛠️ System Maintenance",
		message,
		"system_maintenance",
		clinicID,
	)
}

// NotifyAllUsersAboutSystemUpdate sends system update notification to all users
func NotifyAllUsersAboutSystemUpdate(version string) {
	go SendSystemNotification(
		"🚀 System Update",
		"Dentika has been updated to version "+version+". Please refresh your browser.",
		"system_update",
	)
}
