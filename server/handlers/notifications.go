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
	err := database.DB.Preload("Patient").
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

		// Send reminder via WebSocket
		go SendAppointmentReminder(
			appointment.ID,
			appointment.Patient.FirstName+" "+appointment.Patient.LastName,
			minutesUntil,
		)

		log.Printf("Sent reminder for appointment %d (%d minutes until start)", appointment.ID, minutesUntil)
	}
}
