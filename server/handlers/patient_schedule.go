package handlers

import (
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreatePatientSelfScheduleRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`

	// Appointment preferences
	PreferredDate   string `json:"preferred_date"` // Format: "2006-01-02"
	PreferredTime   string `json:"preferred_time"` // Format: "15:04"
	Symptoms        string `json:"symptoms"`
	AdditionalNotes string `json:"additional_notes"`

	// Branch selection (optional)
	BranchID *uint `json:"branch_id"`
}

type UpdatePatientSelfScheduleRequest struct {
	Status       models.PatientSelfScheduleStatus `json:"status"`
	ReviewNotes  string                           `json:"review_notes"`
	ReviewedByID uint                             `json:"reviewed_by_id"`
}

// CreatePatientSelfSchedule handles patient self-scheduling requests
func CreatePatientSelfSchedule(c *fiber.Ctx) error {
	clinicIdentifier := c.Params("clinicIdentifier")

	// Try to parse as ID first, then as code/name
	var clinic models.Clinic
	if id, err := strconv.ParseUint(clinicIdentifier, 10, 32); err == nil {
		// Look up by ID
		if err := database.DB.First(&clinic, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	} else {
		// Look up by code or name
		if err := database.DB.Where("code = ? OR name = ?", clinicIdentifier, clinicIdentifier).First(&clinic).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	}

	var req CreatePatientSelfScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "First name, last name, email, and phone are required"})
	}

	// Validate branch if provided
	if req.BranchID != nil {
		var branch models.Branch
		if err := database.DB.Where("id = ? AND clinic_id = ?", *req.BranchID, clinic.ID).First(&branch).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid branch for this clinic"})
		}
	}

	// Parse preferred date if provided
	var preferredDate *time.Time
	if req.PreferredDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", req.PreferredDate); err == nil {
			preferredDate = &parsedDate
		}
	}

	// Create the self-schedule request
	scheduleRequest := models.PatientSelfScheduleRequest{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Phone:           req.Phone,
		PreferredDate:   preferredDate,
		PreferredTime:   req.PreferredTime,
		Symptoms:        req.Symptoms,
		AdditionalNotes: req.AdditionalNotes,
		ClinicID:        clinic.ID,
		BranchID:        req.BranchID,
		Status:          models.PatientScheduleStatusPending,
	}

	if err := database.DB.Create(&scheduleRequest).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create self-schedule request"})
	}

	// Reload with relationships
	database.DB.Preload("Clinic").Preload("Branch").First(&scheduleRequest, scheduleRequest.ID)

	return c.Status(201).JSON(scheduleRequest)
}

// GetClinicInfo gets basic clinic information for self-scheduling page
func GetClinicInfo(c *fiber.Ctx) error {
	clinicIdentifier := c.Params("clinicIdentifier")

	// Try to parse as ID first, then as code/name
	var clinic models.Clinic
	if id, err := strconv.ParseUint(clinicIdentifier, 10, 32); err == nil {
		// Look up by ID
		if err := database.DB.First(&clinic, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	} else {
		// Look up by code or name
		if err := database.DB.Where("code = ? OR name = ?", clinicIdentifier, clinicIdentifier).First(&clinic).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	}

	// Load active branches
	database.DB.Preload("Branches", "is_active = ?", true).First(&clinic, clinic.ID)

	return c.JSON(fiber.Map{
		"clinic": clinic,
	})
}

// GetAvailableTimeSlots gets available time slots for a specific date and branch
func GetAvailableTimeSlots(c *fiber.Ctx) error {
	clinicIdentifier := c.Params("clinicIdentifier")
	dateStr := c.Query("date")
	branchIDStr := c.Query("branch_id")

	if dateStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Date is required"})
	}

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format"})
	}
	// Ensure date uses the same location as the server
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	// Get clinic
	var clinic models.Clinic
	if id, err := strconv.ParseUint(clinicIdentifier, 10, 32); err == nil {
		if err := database.DB.First(&clinic, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	} else {
		if err := database.DB.Where("code = ? OR name = ?", clinicIdentifier, clinicIdentifier).First(&clinic).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
		}
	}

	// Get branch if specified
	var branchID *uint
	if branchIDStr != "" {
		if id, err := strconv.ParseUint(branchIDStr, 10, 32); err == nil {
			branchID = &[]uint{uint(id)}[0]
		}
	}

	// Generate available time slots (9 AM to 5 PM in 30-minute increments)
	availableSlots := generateAvailableTimeSlots(date, clinic.ID, branchID)

	return c.JSON(fiber.Map{
		"available_slots": availableSlots,
		"date":            dateStr,
	})
}

func generateAvailableTimeSlots(date time.Time, clinicID uint, branchID *uint) []fiber.Map {
	var availableSlots []fiber.Map

	// Clinic hours: 9 AM to 5 PM
	startHour := 9
	endHour := 17

	// Get existing appointments for this date
	var existingAppointments []models.Appointment

	// Use time range instead of DATE() function to avoid timezone issues
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := database.DB.Where("start_time >= ? AND start_time < ? AND clinic_id = ? AND status IN (?)",
		startOfDay, endOfDay, clinicID,
		[]models.AppointmentStatus{models.StatusScheduled, models.StatusConfirmed})

	// If a specific branch is selected, filter by that branch
	// If no branch is selected, we still want to find all appointments to block those times
	// (since the patient might be assigned to any branch)
	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	query.Find(&existingAppointments)

	// Create a map of occupied time slots
	occupiedSlots := make(map[string]bool)
	for _, appointment := range existingAppointments {
		// Mark the appointment time and any overlapping 30-minute slots
		appointmentStart := appointment.StartTime
		appointmentEnd := appointment.EndTime

		// Generate all 30-minute slots that this appointment occupies
		currentSlot := time.Date(appointmentStart.Year(), appointmentStart.Month(), appointmentStart.Day(),
			appointmentStart.Hour(), appointmentStart.Minute(), 0, 0, time.Local)

		// Round current slot to nearest 30-minute boundary
		minutes := currentSlot.Minute()
		if minutes%30 != 0 {
			// Round up to next 30-minute boundary
			currentSlot = currentSlot.Add(time.Duration(30-minutes%30) * time.Minute)
		}

		// Mark all slots that this appointment occupies
		for currentSlot.Before(appointmentEnd) || currentSlot.Equal(appointmentEnd) {
			if currentSlot.Hour() >= startHour && currentSlot.Hour() < endHour {
				slotKey := currentSlot.Format("15:04")
				occupiedSlots[slotKey] = true
			}
			currentSlot = currentSlot.Add(30 * time.Minute)
		}
	}

	// Generate all possible 30-minute slots
	for hour := startHour; hour < endHour; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.Local)
			slotKey := slotTime.Format("15:04")

			// Check if slot is available
			if !occupiedSlots[slotKey] {
				// Check if slot is in the future (not in the past)
				now := time.Now()
				if slotTime.After(now) || slotTime.Format("2006-01-02") != now.Format("2006-01-02") {
					availableSlots = append(availableSlots, fiber.Map{
						"time":     slotKey,
						"datetime": slotTime.Format("2006-01-02 15:04"),
						"display":  slotTime.Format("3:04 PM"),
					})
				}
			}
		}
	}

	return availableSlots
}
