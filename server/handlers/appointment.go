package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreateAppointmentRequest struct {
	Title               string                   `json:"title"`
	Description         string                   `json:"description"`
	Date                string                   `json:"date"` // Format: "2006-01-02"
	Time                string                   `json:"time"` // Format: "15:04"
	StartTime           time.Time                `json:"-"`    // Calculated from date and time
	EndTime             time.Time                `json:"-"`    // Calculated from start_time and duration
	Duration            int                      `json:"duration"`
	Status              models.AppointmentStatus `json:"status"`
	PatientID           uint                     `json:"patient_id"`
	DoctorID            uint                     `json:"doctor_id"`
	BranchID            uint                     `json:"branch_id"`
	EstimatedCost       *float64                 `json:"estimated_cost,omitempty"`
	PreAppointmentNotes string                   `json:"pre_appointment_notes"`
	// Procedures will be added separately via /appointments/{id}/procedures endpoint
}

// Custom time parsing for JSON
func (t *CreateAppointmentRequest) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to extract date/time
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract date and time
	if dateVal, ok := raw["date"]; ok {
		if dateStr, ok := dateVal.(string); ok {
			t.Date = dateStr
		}
	}
	if timeVal, ok := raw["time"]; ok {
		if timeStr, ok := timeVal.(string); ok {
			t.Time = timeStr
		}
	}

	// Parse date and time to create start_time in Asia/Manila timezone
	if t.Date != "" && t.Time != "" {
		// Load Asia/Manila timezone
		location, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			return err
		}

		// Parse date and time
		dateTimeStr := t.Date + " " + t.Time
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04", dateTimeStr, location)
		if err != nil {
			return err
		}

		t.StartTime = parsedTime

		// Calculate end time based on duration
		if t.Duration > 0 {
			t.EndTime = parsedTime.Add(time.Duration(t.Duration) * time.Minute)
		}
	}

	// Now unmarshal the rest normally
	type Alias CreateAppointmentRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	// Remove date and time from raw to avoid conflict
	delete(raw, "date")
	delete(raw, "time")

	// Marshal back and unmarshal into aux
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	return nil
}

type UpdateAppointmentStatusRequest struct {
	Status               models.AppointmentStatus `json:"status"`
	PostAppointmentNotes string                   `json:"post_appointment_notes"`
	ActualCost           float64                  `json:"actual_cost"`
	NextAppointmentDate  *time.Time               `json:"next_appointment_date"`
}

func GetAppointments(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var appointments []models.Appointment
	query := database.DB.Preload("Patient").Preload("Doctor").Preload("Branch").Preload("Branch.Clinic")

	// Filter by clinic access
	if !user.IsSuperAdmin() {
		// Join with branches to filter by clinic
		query = query.Joins("JOIN branches ON appointments.branch_id = branches.id").
			Where("branches.clinic_id = ?", user.ClinicID)
	}

	// Date filtering
	if dateStr := c.Query("date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			startOfDay := date
			endOfDay := date.Add(24 * time.Hour)
			query = query.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
		}
	}

	// Date range filtering
	if startStr := c.Query("start_date"); startStr != "" {
		if startDate, err := time.Parse("2006-01-02", startStr); err == nil {
			query = query.Where("start_time >= ?", startDate)
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if endDate, err := time.Parse("2006-01-02", endStr); err == nil {
			query = query.Where("start_time <= ?", endDate.Add(24*time.Hour))
		}
	}

	// Status filtering
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Patient filtering
	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Doctor filtering
	if doctorID := c.Query("doctor_id"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Branch filtering
	if branchID := c.Query("branch_id"); branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}

	// Only upcoming appointments
	if c.Query("upcoming") == "true" {
		query = query.Where("start_time > ?", time.Now())
	}

	// Today's appointments
	if c.Query("today") == "true" {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
	}

	// Order by start time
	query = query.Order("start_time ASC")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&appointments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch appointments"})
	}

	return c.JSON(fiber.Map{
		"appointments": appointments,
		"page":         page,
		"limit":        limit,
	})
}

func GetAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	var appointment models.Appointment
	query := database.DB.Preload("Patient").Preload("Doctor").Preload("Branch").Preload("Branch.Clinic").
		Preload("Procedures").Preload("Procedures.ProcedureTemplate").Preload("Procedures.PerformedBy").
		Preload("Diagnoses").Preload("Diagnoses.DiagnosisTemplate").Preload("Diagnoses.DiagnosedBy")

	if err := query.First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != appointment.Branch.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	return c.JSON(appointment)
}

func CreateAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req CreateAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		// Log the raw request body for debugging
		body := c.Body()
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request",
			"details": err.Error(),
			"body":    string(body),
		})
	}

	// Validate required fields
	if req.PatientID == 0 || req.DoctorID == 0 || req.BranchID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Patient, doctor, and branch are required"})
	}

	if req.Date == "" || req.Time == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Date and time are required"})
	}

	if req.StartTime.IsZero() {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date or time format"})
	}

	// Validate patient belongs to accessible clinic
	var patient models.Patient
	if err := database.DB.First(&patient, req.PatientID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
	}

	// Validate doctor belongs to accessible clinic
	var doctor models.User
	if err := database.DB.First(&doctor, req.DoctorID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Doctor not found"})
	}

	// Validate branch belongs to accessible clinic
	var branch models.Branch
	if err := database.DB.First(&branch, req.BranchID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Branch not found"})
	}

	// Check access to clinic
	if !user.IsSuperAdmin() {
		if user.ClinicID != patient.ClinicID ||
			user.ClinicID != branch.ClinicID {
			return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
		}

		// Check if doctor belongs to same clinic (skip for super admin doctors)
		if !doctor.IsSuperAdmin() && doctor.ClinicID != user.ClinicID {
			return c.Status(400).JSON(fiber.Map{"error": "Doctor does not belong to the same clinic"})
		}
	}

	// Calculate end time if not provided
	if req.EndTime.IsZero() {
		if req.Duration > 0 {
			req.EndTime = req.StartTime.Add(time.Duration(req.Duration) * time.Minute)
		} else {
			req.EndTime = req.StartTime.Add(30 * time.Minute) // Default 30 minutes
			req.Duration = 30
		}
	} else if req.Duration == 0 {
		req.Duration = int(req.EndTime.Sub(req.StartTime).Minutes())
	}

	// Check for conflicting appointments with same doctor
	var conflictingAppointments []models.Appointment
	if err := database.DB.Where("doctor_id = ? AND status IN (?, ?) AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		req.DoctorID, models.StatusScheduled, models.StatusConfirmed,
		req.StartTime, req.StartTime, req.EndTime, req.EndTime).Find(&conflictingAppointments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check for conflicts"})
	}

	if len(conflictingAppointments) > 0 {
		return c.Status(409).JSON(fiber.Map{
			"error":     "Doctor has conflicting appointment at this time",
			"conflicts": conflictingAppointments,
		})
	}

	// Set default estimated cost if not provided
	estimatedCost := 0.0
	if req.EstimatedCost != nil {
		estimatedCost = *req.EstimatedCost
	}

	appointment := models.Appointment{
		Title:               req.Title,
		Description:         req.Description,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
		Duration:            req.Duration,
		Status:              req.Status,
		PatientID:           req.PatientID,
		DoctorID:            req.DoctorID,
		BranchID:            req.BranchID,
		ClinicID:            branch.ClinicID,
		EstimatedCost:       estimatedCost,
		PreAppointmentNotes: req.PreAppointmentNotes,
	}

	// Set default status if not provided
	if appointment.Status == "" {
		appointment.Status = models.StatusScheduled
	}

	if err := database.DB.Create(&appointment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create appointment"})
	}

	// Reload with relationships
	database.DB.Preload("Patient").Preload("Doctor").Preload("Branch").First(&appointment, appointment.ID)

	// Send WebSocket notification for new appointment (clinic-scoped)
	go SendAppointmentUpdate(appointment.ID, appointment.Patient.FirstName+" "+appointment.Patient.LastName, "scheduled", appointment.Branch.ClinicID)

	// Send clinic-wide notification
	go SendClinicNotification(
		"New Appointment Scheduled",
		"Appointment scheduled for "+appointment.Patient.FirstName+" "+appointment.Patient.LastName+" with Dr. "+appointment.Doctor.FirstName+" "+appointment.Doctor.LastName,
		"appointment_scheduled",
		appointment.Branch.ClinicID,
	)

	return c.Status(201).JSON(appointment)
}

func UpdateAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != appointment.Branch.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var req CreateAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update appointment fields
	if req.Title != "" {
		appointment.Title = req.Title
	}
	if req.Description != "" {
		appointment.Description = req.Description
	}
	if !req.StartTime.IsZero() {
		appointment.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		appointment.EndTime = req.EndTime
	}
	if req.Duration > 0 {
		appointment.Duration = req.Duration
	}
	if req.Status != "" {
		appointment.Status = req.Status
	}

	// Update foreign key fields
	if req.PatientID > 0 {
		appointment.PatientID = req.PatientID
	}
	if req.DoctorID > 0 {
		appointment.DoctorID = req.DoctorID
	}
	if req.BranchID > 0 {
		appointment.BranchID = req.BranchID
		// Update clinic ID based on branch
		var branch models.Branch
		if err := database.DB.First(&branch, req.BranchID).Error; err == nil {
			appointment.ClinicID = branch.ClinicID
		}
	}

	if req.EstimatedCost != nil {
		appointment.EstimatedCost = *req.EstimatedCost
	}
	if req.PreAppointmentNotes != "" {
		appointment.PreAppointmentNotes = req.PreAppointmentNotes
	}

	// Recalculate duration if times changed
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		appointment.Duration = int(appointment.EndTime.Sub(appointment.StartTime).Minutes())
	}

	if err := database.DB.Save(&appointment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update appointment"})
	}

	// Reload with relationships
	database.DB.Preload("Patient").Preload("Doctor").Preload("Branch").First(&appointment, appointment.ID)

	// Send WebSocket notification for appointment update (clinic-scoped)
	go SendAppointmentUpdate(appointment.ID, appointment.Patient.FirstName+" "+appointment.Patient.LastName, "updated", appointment.Branch.ClinicID)

	return c.JSON(appointment)
}

func UpdateAppointmentStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != appointment.Branch.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var req UpdateAppointmentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	appointment.Status = req.Status
	if req.PostAppointmentNotes != "" {
		appointment.PostAppointmentNotes = req.PostAppointmentNotes
	}
	if req.ActualCost > 0 {
		appointment.ActualCost = req.ActualCost
	}
	if req.NextAppointmentDate != nil {
		appointment.NextAppointmentDate = req.NextAppointmentDate
	}

	if err := database.DB.Save(&appointment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update appointment status"})
	}

	// Reload with relationships for notification
	database.DB.Preload("Patient").Preload("Branch").First(&appointment, appointment.ID)

	// Send WebSocket notification for status update (clinic-scoped)
	go SendAppointmentUpdate(appointment.ID, appointment.Patient.FirstName+" "+appointment.Patient.LastName, string(appointment.Status), appointment.Branch.ClinicID)

	// Send clinic-wide notification for status changes
	statusMessages := map[string]string{
		"confirmed": "Appointment confirmed for " + appointment.Patient.FirstName + " " + appointment.Patient.LastName,
		"completed": "Appointment completed for " + appointment.Patient.FirstName + " " + appointment.Patient.LastName,
		"cancelled": "Appointment cancelled for " + appointment.Patient.FirstName + " " + appointment.Patient.LastName,
		"no_show":   "Patient no-show for " + appointment.Patient.FirstName + " " + appointment.Patient.LastName,
	}

	if message, exists := statusMessages[string(appointment.Status)]; exists {
		go SendClinicNotification(
			"Appointment Status Update",
			message,
			"appointment_"+string(appointment.Status),
			appointment.Branch.ClinicID,
		)
	}

	return c.JSON(appointment)
}

func MarkPatientArrived(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != appointment.Branch.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	appointment.MarkPatientArrived()

	if err := database.DB.Save(&appointment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to mark patient as arrived"})
	}

	return c.JSON(fiber.Map{"message": "Patient marked as arrived", "is_late": appointment.IsLate})
}

func CheckAppointmentAvailability(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req struct {
		DoctorID  uint   `json:"doctor_id"`
		Date      string `json:"date"`
		Time      string `json:"time"`
		Duration  int    `json:"duration"`
		BranchID  uint   `json:"branch_id,omitempty"`
		ExcludeID uint   `json:"exclude_id,omitempty"` // For editing existing appointments
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate required fields
	if req.DoctorID == 0 || req.Date == "" || req.Time == "" || req.Duration == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Doctor ID, date, time, and duration are required"})
	}

	// Parse date and time to create start_time in Asia/Manila timezone
	location, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load timezone"})
	}

	dateTimeStr := req.Date + " " + req.Time
	startTime, err := time.ParseInLocation("2006-01-02 15:04", dateTimeStr, location)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date or time format"})
	}

	endTime := startTime.Add(time.Duration(req.Duration) * time.Minute)

	// Check for conflicting appointments
	query := database.DB.Where("doctor_id = ? AND status IN (?, ?) AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		req.DoctorID, models.StatusScheduled, models.StatusConfirmed,
		startTime, startTime, endTime, endTime)

	// Exclude the appointment being edited (if editing)
	if req.ExcludeID > 0 {
		query = query.Where("id != ?", req.ExcludeID)
	}

	// Filter by clinic access if not super admin
	if !user.IsSuperAdmin() {
		query = query.Joins("JOIN branches ON appointments.branch_id = branches.id").
			Where("branches.clinic_id = ?", user.ClinicID)
	}

	var conflictingAppointments []models.Appointment
	if err := query.Preload("Patient").Find(&conflictingAppointments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check for conflicts"})
	}

	// Return availability status
	if len(conflictingAppointments) > 0 {
		return c.JSON(fiber.Map{
			"available": false,
			"conflicts": conflictingAppointments,
			"message":   "Doctor has conflicting appointment(s) at this time",
		})
	}

	return c.JSON(fiber.Map{
		"available": true,
		"message":   "Time slot is available",
	})
}

func GetUpcomingAppointments(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var appointments []models.Appointment
	query := database.DB.Preload("Patient").Preload("Doctor").Preload("Branch").
		Where("start_time > ? AND status IN (?)", time.Now(), []models.AppointmentStatus{
			models.StatusScheduled, models.StatusConfirmed,
		})

	// Filter by clinic access
	if !user.IsSuperAdmin() {
		query = query.Joins("JOIN branches ON appointments.branch_id = branches.id").
			Where("branches.clinic_id = ?", user.ClinicID)
	}

	// Get appointments for next 30 minutes for countdown
	thirtyMinutesFromNow := time.Now().Add(30 * time.Minute)
	query = query.Where("start_time <= ?", thirtyMinutesFromNow)

	query = query.Order("start_time ASC").Limit(10)

	if err := query.Find(&appointments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch upcoming appointments"})
	}

	// Add countdown information
	var appointmentsWithCountdown []fiber.Map
	for _, apt := range appointments {
		timeUntil := apt.GetTimeUntilAppointment()
		appointmentsWithCountdown = append(appointmentsWithCountdown, fiber.Map{
			"appointment":           apt,
			"time_until_minutes":    int(timeUntil.Minutes()),
			"should_show_countdown": apt.ShouldShowCountdown(),
		})
	}

	return c.JSON(appointmentsWithCountdown)
}
