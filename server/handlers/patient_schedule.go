package handlers

import (
	"encoding/json"
	"strconv"
	"strings"
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

	// Get the operating hours for this date
	operatingHours := getBranchOperatingHours(date, clinicID, branchID)
	if len(operatingHours) == 0 {
		// No operating hours available for this date
		return availableSlots
	}

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
			// Mark the slot as occupied regardless of operating hours
			// We'll filter by operating hours when generating available slots
			slotKey := currentSlot.Format("15:04")
			occupiedSlots[slotKey] = true
			currentSlot = currentSlot.Add(30 * time.Minute)
		}
	}

	// Generate all possible 30-minute time slots for each operating period
	now := time.Now()
	for _, period := range operatingHours {
		// Parse start and end times
		startTime, err := time.Parse("15:04", period.Start)
		if err != nil {
			continue // Skip invalid time format
		}
		endTime, err := time.Parse("15:04", period.End)
		if err != nil {
			continue // Skip invalid time format
		}

		startHour := startTime.Hour()
		startMinute := startTime.Minute()
		endHour := endTime.Hour()
		endMinute := endTime.Minute()

		// Generate slots for this period
		currentSlot := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMinute, 0, 0, time.Local)
		endSlot := time.Date(date.Year(), date.Month(), date.Day(), endHour, endMinute, 0, 0, time.Local)

		// Round start time to nearest 30-minute boundary (up)
		if currentSlot.Minute()%30 != 0 {
			minutes := currentSlot.Minute()
			currentSlot = currentSlot.Add(time.Duration(30-minutes%30) * time.Minute)
		}

		for currentSlot.Before(endSlot) {
			slotKey := currentSlot.Format("15:04")

			// Skip slots that are occupied by existing appointments
			if !occupiedSlots[slotKey] {
				// Only include slots that are in the future (if today) or any slot for future dates
				if currentSlot.After(now) || currentSlot.Format("2006-01-02") != now.Format("2006-01-02") {
					availableSlots = append(availableSlots, fiber.Map{
						"time":     slotKey,
						"datetime": currentSlot.Format("2006-01-02 15:04"),
						"display":  currentSlot.Format("3:04 PM"),
					})
				}
			}
			currentSlot = currentSlot.Add(30 * time.Minute)
		}
	}

	return availableSlots
}

// SchedulePeriod represents a time period (e.g., 09:00-13:00)
type SchedulePeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// DaySchedule represents the schedule for a specific day
type DaySchedule struct {
	Enabled bool             `json:"enabled"`
	Periods []SchedulePeriod `json:"periods"`
}

// BranchSchedule represents the complete schedule configuration
type BranchSchedule struct {
	Monday    DaySchedule `json:"monday"`
	Tuesday   DaySchedule `json:"tuesday"`
	Wednesday DaySchedule `json:"wednesday"`
	Thursday  DaySchedule `json:"thursday"`
	Friday    DaySchedule `json:"friday"`
	Saturday  DaySchedule `json:"saturday"`
	Sunday    DaySchedule `json:"sunday"`
	Holidays  struct {
		Enabled bool `json:"enabled"`
	} `json:"holidays"`
	Timezone string `json:"timezone"`
}

// getBranchOperatingHours returns the operating hours for a given date and branch
func getBranchOperatingHours(date time.Time, clinicID uint, branchID *uint) []SchedulePeriod {
	// Get branches
	var branches []models.Branch
	query := database.DB.Where("clinic_id = ? AND is_active = ?", clinicID, true)

	if branchID != nil {
		query = query.Where("id = ?", *branchID)
	}

	if err := query.Find(&branches).Error; err != nil {
		return []SchedulePeriod{}
	}

	if len(branches) == 0 {
		return []SchedulePeriod{}
	}

	// Use the first branch (or the specified branch)
	branch := branches[0]

	// Check if the branch is closed today (emergency override)
	if branch.IsClosedToday {
		return []SchedulePeriod{}
	}

	// Parse the schedule JSON
	if branch.Schedule == "" {
		// Fallback to default hours if no schedule is set
		return []SchedulePeriod{
			{Start: "09:00", End: "17:00"},
		}
	}

	var schedule BranchSchedule
	if err := json.Unmarshal([]byte(branch.Schedule), &schedule); err != nil {
		// Fallback to default hours if JSON is invalid
		return []SchedulePeriod{
			{Start: "09:00", End: "17:00"},
		}
	}

	// Get the day of the week
	dayName := strings.ToLower(date.Weekday().String())

	var daySchedule DaySchedule
	switch dayName {
	case "monday":
		daySchedule = schedule.Monday
	case "tuesday":
		daySchedule = schedule.Tuesday
	case "wednesday":
		daySchedule = schedule.Wednesday
	case "thursday":
		daySchedule = schedule.Thursday
	case "friday":
		daySchedule = schedule.Friday
	case "saturday":
		daySchedule = schedule.Saturday
	case "sunday":
		daySchedule = schedule.Sunday
	default:
		return []SchedulePeriod{}
	}

	// Check if the day is enabled
	if !daySchedule.Enabled {
		return []SchedulePeriod{}
	}

	return daySchedule.Periods
}
