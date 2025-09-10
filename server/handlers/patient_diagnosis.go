package handlers

import (
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreatePatientDiagnosisRequest struct {
	DiagnosisTemplateID uint   `json:"diagnosis_template_id" validate:"required"`
	ToothNumber         string `json:"tooth_number"`
	Surface             string `json:"surface"`
	Notes               string `json:"notes"`
	Severity            string `json:"severity"`
	AppointmentID       *uint  `json:"appointment_id"`
}

type UpdatePatientDiagnosisRequest struct {
	ToothNumber   string `json:"tooth_number"`
	Surface       string `json:"surface"`
	Notes         string `json:"notes"`
	Severity      string `json:"severity"`
	Status        string `json:"status"`
	AppointmentID *uint  `json:"appointment_id"`
}

type CreatePatientTreatmentPlanRequest struct {
	Title             string                          `json:"title" validate:"required"`
	Description       string                          `json:"description"`
	Priority          string                          `json:"priority"`
	DiagnosisID       *uint                           `json:"diagnosis_id"`
	EstimatedCost     float64                         `json:"estimated_cost"`
	EstimatedVisits   int                             `json:"estimated_visits"`
	EstimatedDuration int                             `json:"estimated_duration"`
	StartDate         *time.Time                      `json:"start_date"`
	TargetCompletion  *time.Time                      `json:"target_completion"`
	Procedures        []TreatmentPlanProcedureRequest `json:"procedures"`
}

type TreatmentPlanProcedureRequest struct {
	ProcedureTemplateID uint    `json:"procedure_template_id" validate:"required"`
	ToothNumber         string  `json:"tooth_number"`
	Surface             string  `json:"surface"`
	Notes               string  `json:"notes"`
	EstimatedCost       float64 `json:"estimated_cost"`
	Sequence            int     `json:"sequence"`
}

type UpdatePatientTreatmentPlanRequest struct {
	Title             string                          `json:"title"`
	Description       string                          `json:"description"`
	Priority          string                          `json:"priority"`
	Status            string                          `json:"status"`
	DiagnosisID       *uint                           `json:"diagnosis_id"`
	EstimatedCost     float64                         `json:"estimated_cost"`
	EstimatedVisits   int                             `json:"estimated_visits"`
	EstimatedDuration int                             `json:"estimated_duration"`
	StartDate         *time.Time                      `json:"start_date"`
	TargetCompletion  *time.Time                      `json:"target_completion"`
	CompletedVisits   int                             `json:"completed_visits"`
	ActualCost        float64                         `json:"actual_cost"`
	Procedures        []TreatmentPlanProcedureRequest `json:"procedures"`
}

// GetPatientDiagnoses - Get all diagnoses for a specific patient
func GetPatientDiagnoses(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var diagnoses []models.PatientDiagnosis
	query := database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").Preload("Appointment")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	query = query.Where("patient_id = ?", patientID)

	// Add status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Add sorting
	query = query.Order("diagnosed_at DESC")

	if err := query.Find(&diagnoses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch patient diagnoses"})
	}

	return c.JSON(fiber.Map{
		"diagnoses": diagnoses,
	})
}

// GetPatientDiagnosis - Get single patient diagnosis
func GetPatientDiagnosis(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	diagnosisID, err := strconv.ParseUint(c.Params("diagnosisId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid diagnosis ID"})
	}

	var diagnosis models.PatientDiagnosis
	query := database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").Preload("Appointment")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&diagnosis, diagnosisID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient diagnosis not found"})
	}

	return c.JSON(diagnosis)
}

// CreatePatientDiagnosis - Create new patient diagnosis
func CreatePatientDiagnosis(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var req CreatePatientDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.DiagnosisTemplateID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Diagnosis template ID is required"})
	}

	// Get clinic ID
	var clinicID uint
	if user.IsSuperAdmin() {
		// For super admin, get clinic from patient
		var patient models.Patient
		if err := database.DB.First(&patient, patientID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
		}
		clinicID = patient.ClinicID
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	// Verify diagnosis template exists
	var template models.DiagnosisTemplate
	if err := database.DB.First(&template, req.DiagnosisTemplateID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid diagnosis template ID"})
	}

	diagnosis := models.PatientDiagnosis{
		DiagnosisTemplateID: req.DiagnosisTemplateID,
		PatientID:           uint(patientID),
		ClinicID:            clinicID,
		ToothNumber:         req.ToothNumber,
		Surface:             req.Surface,
		Notes:               req.Notes,
		Severity:            req.Severity,
		AppointmentID:       req.AppointmentID,
		DiagnosedByID:       user.ID,
		DiagnosedAt:         time.Now(),
	}

	if err := database.DB.Create(&diagnosis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create patient diagnosis"})
	}

	// Preload relationships for response
	database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").Preload("Appointment").First(&diagnosis, diagnosis.ID)

	return c.Status(201).JSON(diagnosis)
}

// UpdatePatientDiagnosis - Update existing patient diagnosis
func UpdatePatientDiagnosis(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	diagnosisID, err := strconv.ParseUint(c.Params("diagnosisId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid diagnosis ID"})
	}

	var req UpdatePatientDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var diagnosis models.PatientDiagnosis
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&diagnosis, diagnosisID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient diagnosis not found"})
	}

	// Update fields
	if req.ToothNumber != "" {
		diagnosis.ToothNumber = req.ToothNumber
	}
	if req.Surface != "" {
		diagnosis.Surface = req.Surface
	}
	if req.Notes != "" {
		diagnosis.Notes = req.Notes
	}
	if req.Severity != "" {
		diagnosis.Severity = req.Severity
	}
	if req.Status != "" {
		diagnosis.Status = req.Status
		if req.Status == "resolved" {
			now := time.Now()
			diagnosis.ResolvedAt = &now
		}
	}
	if req.AppointmentID != nil {
		diagnosis.AppointmentID = req.AppointmentID
	}

	diagnosis.UpdatedAt = time.Now()

	if err := database.DB.Save(&diagnosis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update patient diagnosis"})
	}

	return c.JSON(diagnosis)
}

// DeletePatientDiagnosis - Delete patient diagnosis
func DeletePatientDiagnosis(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	diagnosisID, err := strconv.ParseUint(c.Params("diagnosisId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid diagnosis ID"})
	}

	var diagnosis models.PatientDiagnosis
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&diagnosis, diagnosisID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient diagnosis not found"})
	}

	if err := database.DB.Delete(&diagnosis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete patient diagnosis"})
	}

	return c.JSON(fiber.Map{"message": "Patient diagnosis deleted successfully"})
}

// GetPatientTreatmentPlans - Get all treatment plans for a specific patient
func GetPatientTreatmentPlans(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var treatmentPlans []models.PatientTreatmentPlan
	query := database.DB.Preload("CreatedBy").Preload("Diagnosis").Preload("Diagnosis.DiagnosisTemplate")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	query = query.Where("patient_id = ?", patientID)

	// Add status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Add sorting
	query = query.Order("created_at DESC")

	if err := query.Find(&treatmentPlans).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch patient treatment plans"})
	}

	return c.JSON(fiber.Map{
		"treatment_plans": treatmentPlans,
	})
}

// GetPatientTreatmentPlan - Get single patient treatment plan
func GetPatientTreatmentPlan(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	treatmentPlanID, err := strconv.ParseUint(c.Params("treatmentPlanId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid treatment plan ID"})
	}

	var treatmentPlan models.PatientTreatmentPlan
	query := database.DB.Preload("CreatedBy").Preload("Diagnosis").Preload("Diagnosis.DiagnosisTemplate")

	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&treatmentPlan, treatmentPlanID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient treatment plan not found"})
	}

	return c.JSON(treatmentPlan)
}

// CreatePatientTreatmentPlan - Create new patient treatment plan
func CreatePatientTreatmentPlan(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var req CreatePatientTreatmentPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Treatment plan title is required"})
	}

	// Get clinic ID
	var clinicID uint
	if user.IsSuperAdmin() {
		// For super admin, get clinic from patient
		var patient models.Patient
		if err := database.DB.First(&patient, patientID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
		}
		clinicID = patient.ClinicID
	} else {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		clinicID = *user.ClinicID
	}

	treatmentPlan := models.PatientTreatmentPlan{
		PatientID:         uint(patientID),
		ClinicID:          clinicID,
		Title:             req.Title,
		Description:       req.Description,
		Priority:          req.Priority,
		DiagnosisID:       req.DiagnosisID,
		EstimatedCost:     req.EstimatedCost,
		EstimatedVisits:   req.EstimatedVisits,
		EstimatedDuration: req.EstimatedDuration,
		StartDate:         req.StartDate,
		TargetCompletion:  req.TargetCompletion,
		CreatedByID:       user.ID,
	}

	if err := database.DB.Create(&treatmentPlan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create patient treatment plan"})
	}

	// Create treatment plan procedures if provided
	if len(req.Procedures) > 0 {
		for _, procReq := range req.Procedures {
			procedure := models.TreatmentPlanProcedure{
				TreatmentPlanID:     treatmentPlan.ID,
				ProcedureTemplateID: procReq.ProcedureTemplateID,
				ToothNumber:         procReq.ToothNumber,
				Surface:             procReq.Surface,
				Notes:               procReq.Notes,
				EstimatedCost:       procReq.EstimatedCost,
				Sequence:            procReq.Sequence,
			}
			database.DB.Create(&procedure)
		}
	}

	// Preload relationships for response
	database.DB.Preload("CreatedBy").Preload("Diagnosis").Preload("Diagnosis.DiagnosisTemplate").First(&treatmentPlan, treatmentPlan.ID)

	return c.Status(201).JSON(treatmentPlan)
}

// UpdatePatientTreatmentPlan - Update existing patient treatment plan
func UpdatePatientTreatmentPlan(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	treatmentPlanID, err := strconv.ParseUint(c.Params("treatmentPlanId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid treatment plan ID"})
	}

	var req UpdatePatientTreatmentPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var treatmentPlan models.PatientTreatmentPlan
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&treatmentPlan, treatmentPlanID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient treatment plan not found"})
	}

	// Update fields
	if req.Title != "" {
		treatmentPlan.Title = req.Title
	}
	if req.Description != "" {
		treatmentPlan.Description = req.Description
	}
	if req.Priority != "" {
		treatmentPlan.Priority = req.Priority
	}
	if req.Status != "" {
		treatmentPlan.Status = req.Status
		if req.Status == "completed" {
			now := time.Now()
			treatmentPlan.CompletedAt = &now
		}
	}
	if req.DiagnosisID != nil {
		treatmentPlan.DiagnosisID = req.DiagnosisID
	}
	if req.EstimatedCost > 0 {
		treatmentPlan.EstimatedCost = req.EstimatedCost
	}
	if req.EstimatedVisits > 0 {
		treatmentPlan.EstimatedVisits = req.EstimatedVisits
	}
	if req.EstimatedDuration > 0 {
		treatmentPlan.EstimatedDuration = req.EstimatedDuration
	}
	if req.StartDate != nil {
		treatmentPlan.StartDate = req.StartDate
	}
	if req.TargetCompletion != nil {
		treatmentPlan.TargetCompletion = req.TargetCompletion
	}
	if req.CompletedVisits >= 0 {
		treatmentPlan.CompletedVisits = req.CompletedVisits
	}
	if req.ActualCost >= 0 {
		treatmentPlan.ActualCost = req.ActualCost
	}

	treatmentPlan.UpdatedAt = time.Now()

	if err := database.DB.Save(&treatmentPlan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update patient treatment plan"})
	}

	return c.JSON(treatmentPlan)
}

// DeletePatientTreatmentPlan - Delete patient treatment plan
func DeletePatientTreatmentPlan(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patientId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	treatmentPlanID, err := strconv.ParseUint(c.Params("treatmentPlanId"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid treatment plan ID"})
	}

	var treatmentPlan models.PatientTreatmentPlan
	query := database.DB
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}
		query = query.Where("clinic_id = ?", *user.ClinicID)
	}

	if err := query.Where("patient_id = ?", patientID).First(&treatmentPlan, treatmentPlanID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient treatment plan not found"})
	}

	if err := database.DB.Delete(&treatmentPlan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete patient treatment plan"})
	}

	return c.JSON(fiber.Map{"message": "Patient treatment plan deleted successfully"})
}
