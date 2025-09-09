package handlers

import (
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

// Procedure Template Handlers
func GetProcedureTemplates(c *fiber.Ctx) error {
	var templates []models.ProcedureTemplate
	query := database.DB.Where("is_active = ?", true)

	// Filter by category if provided
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Search by name or code
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("category ASC, name ASC").Find(&templates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch procedure templates"})
	}

	return c.JSON(templates)
}

func CreateProcedureTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Only super admin and clinic owners can create templates
	if !user.IsSuperAdmin() && !user.HasRole(models.ClinicOwner) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	var req models.ProcedureTemplate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Code == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Code and name are required"})
	}

	// Check for duplicate code
	var existing models.ProcedureTemplate
	if err := database.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Procedure code already exists"})
	}

	req.IsActive = true

	if err := database.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create procedure template"})
	}

	return c.Status(201).JSON(req)
}

// Diagnosis Template Handlers
func GetDiagnosisTemplates(c *fiber.Ctx) error {
	var templates []models.DiagnosisTemplate
	query := database.DB.Where("is_active = ?", true)

	// Filter by category if provided
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Filter by severity if provided
	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	// Search by name or code
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("category ASC, name ASC").Find(&templates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch diagnosis templates"})
	}

	return c.JSON(templates)
}

func CreateDiagnosisTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Only super admin and clinic owners can create templates
	if !user.IsSuperAdmin() && !user.HasRole(models.ClinicOwner) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	var req models.DiagnosisTemplate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Code == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Code and name are required"})
	}

	// Check for duplicate code
	var existing models.DiagnosisTemplate
	if err := database.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Diagnosis code already exists"})
	}

	req.IsActive = true

	if err := database.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create diagnosis template"})
	}

	return c.Status(201).JSON(req)
}

// Appointment Procedure Handlers
func GetAppointmentProcedures(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("appointment_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	// Verify appointment access
	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var procedures []models.AppointmentProcedure
	if err := database.DB.Preload("ProcedureTemplate").Preload("PerformedBy").
		Where("appointment_id = ?", appointmentID).Find(&procedures).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch procedures"})
	}

	return c.JSON(procedures)
}

func AddProcedureToAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("appointment_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	// Verify appointment access
	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can add procedures
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can add procedures"})
	}

	type AddProcedureRequest struct {
		ProcedureTemplateID uint    `json:"procedure_template_id"`
		ToothNumber         string  `json:"tooth_number"`
		Surface             string  `json:"surface"`
		Notes               string  `json:"notes"`
		Cost                float64 `json:"cost"`
		Status              string  `json:"status"`
	}

	var req AddProcedureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.ProcedureTemplateID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Procedure template is required"})
	}

	// Verify procedure template exists
	var template models.ProcedureTemplate
	if err := database.DB.First(&template, req.ProcedureTemplateID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Procedure template not found"})
	}

	procedure := models.AppointmentProcedure{
		ProcedureTemplateID: req.ProcedureTemplateID,
		AppointmentID:       uint(appointmentID),
		ToothNumber:         req.ToothNumber,
		Surface:             req.Surface,
		Notes:               req.Notes,
		Cost:                req.Cost,
		Status:              req.Status,
		PerformedByID:       user.ID,
	}

	// Set default cost if not provided
	if procedure.Cost == 0 {
		procedure.Cost = template.DefaultCost
	}

	// Set default status if not provided
	if procedure.Status == "" {
		procedure.Status = "planned"
	}

	if err := database.DB.Create(&procedure).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add procedure"})
	}

	// Reload with relationships
	database.DB.Preload("ProcedureTemplate").Preload("PerformedBy").First(&procedure, procedure.ID)

	return c.Status(201).JSON(procedure)
}

func UpdateAppointmentProcedure(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	procedureID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid procedure ID"})
	}

	var procedure models.AppointmentProcedure
	if err := database.DB.Preload("Appointment").Preload("Appointment.Branch").First(&procedure, procedureID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Procedure not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != procedure.Appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can modify procedures
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can modify procedures"})
	}

	type UpdateProcedureRequest struct {
		ToothNumber string     `json:"tooth_number"`
		Surface     string     `json:"surface"`
		Notes       string     `json:"notes"`
		Cost        float64    `json:"cost"`
		Status      string     `json:"status"`
		StartTime   *time.Time `json:"start_time"`
		EndTime     *time.Time `json:"end_time"`
	}

	var req UpdateProcedureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update fields
	if req.ToothNumber != "" {
		procedure.ToothNumber = req.ToothNumber
	}
	if req.Surface != "" {
		procedure.Surface = req.Surface
	}
	if req.Notes != "" {
		procedure.Notes = req.Notes
	}
	if req.Cost > 0 {
		procedure.Cost = req.Cost
	}
	if req.Status != "" {
		procedure.Status = req.Status
	}
	if req.StartTime != nil {
		procedure.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		procedure.EndTime = req.EndTime
	}

	if err := database.DB.Save(&procedure).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update procedure"})
	}

	// Reload with relationships
	database.DB.Preload("ProcedureTemplate").Preload("PerformedBy").First(&procedure, procedure.ID)

	return c.JSON(procedure)
}

// Appointment Diagnosis Handlers
func GetAppointmentDiagnoses(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("appointment_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	// Verify appointment access
	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var diagnoses []models.AppointmentDiagnosis
	if err := database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").
		Where("appointment_id = ?", appointmentID).Find(&diagnoses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch diagnoses"})
	}

	return c.JSON(diagnoses)
}

func AddDiagnosisToAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	appointmentID, err := strconv.ParseUint(c.Params("appointment_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid appointment ID"})
	}

	// Verify appointment access
	var appointment models.Appointment
	if err := database.DB.Preload("Branch").First(&appointment, appointmentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Appointment not found"})
	}

	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can add diagnoses
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can add diagnoses"})
	}

	type AddDiagnosisRequest struct {
		DiagnosisTemplateID uint   `json:"diagnosis_template_id"`
		ToothNumber         string `json:"tooth_number"`
		Surface             string `json:"surface"`
		Notes               string `json:"notes"`
		Severity            string `json:"severity"`
		TreatmentPlan       string `json:"treatment_plan"`
	}

	var req AddDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.DiagnosisTemplateID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Diagnosis template is required"})
	}

	// Verify diagnosis template exists
	var template models.DiagnosisTemplate
	if err := database.DB.First(&template, req.DiagnosisTemplateID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Diagnosis template not found"})
	}

	diagnosis := models.AppointmentDiagnosis{
		DiagnosisTemplateID: req.DiagnosisTemplateID,
		AppointmentID:       uint(appointmentID),
		ToothNumber:         req.ToothNumber,
		Surface:             req.Surface,
		Notes:               req.Notes,
		Severity:            req.Severity,
		TreatmentPlan:       req.TreatmentPlan,
		DiagnosedByID:       user.ID,
	}

	// Set default severity if not provided
	if diagnosis.Severity == "" {
		diagnosis.Severity = template.Severity
	}

	if err := database.DB.Create(&diagnosis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add diagnosis"})
	}

	// Reload with relationships
	database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").First(&diagnosis, diagnosis.ID)

	return c.Status(201).JSON(diagnosis)
}

func UpdateAppointmentDiagnosis(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	diagnosisID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid diagnosis ID"})
	}

	var diagnosis models.AppointmentDiagnosis
	if err := database.DB.Preload("Appointment").Preload("Appointment.Branch").First(&diagnosis, diagnosisID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Diagnosis not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && (user.ClinicID == nil || *user.ClinicID != diagnosis.Appointment.Branch.ClinicID) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can modify diagnoses
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can modify diagnoses"})
	}

	type UpdateDiagnosisRequest struct {
		ToothNumber   string `json:"tooth_number"`
		Surface       string `json:"surface"`
		Notes         string `json:"notes"`
		Severity      string `json:"severity"`
		TreatmentPlan string `json:"treatment_plan"`
	}

	var req UpdateDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
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
	if req.TreatmentPlan != "" {
		diagnosis.TreatmentPlan = req.TreatmentPlan
	}

	if err := database.DB.Save(&diagnosis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update diagnosis"})
	}

	// Reload with relationships
	database.DB.Preload("DiagnosisTemplate").Preload("DiagnosedBy").First(&diagnosis, diagnosis.ID)

	return c.JSON(diagnosis)
}
