package handlers

import (
	"strings"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

// Consent Template Handlers
func GetConsentTemplates(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	var templates []models.ConsentTemplate

	query := database.DB.Where("is_active = ?", true)

	// Filter by clinic for non-super-admin users
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch consent templates"})
	}

	return c.JSON(templates)
}

func CreateConsentTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Only super admin and admins can create templates
	if !user.IsSuperAdmin() && !user.HasRole(models.Admin) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	var req models.ConsentTemplate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Code == "" || req.Name == "" || req.Content == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Code, name, and content are required"})
	}

	// Set clinic ID for non-super-admin users
	if !user.IsSuperAdmin() {
		req.ClinicID = user.ClinicID
	}

	// Check for duplicate code within the same clinic
	var existing models.ConsentTemplate
	query := database.DB.Where("code = ?", req.Code)
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", req.ClinicID)
	}
	if err := query.First(&existing).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Consent template code already exists"})
	}

	req.IsActive = true

	if err := database.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create consent template"})
	}

	return c.Status(201).JSON(req)
}

func GetConsentTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	var template models.ConsentTemplate
	query := database.DB.Where("id = ? AND is_active = ?", id, true)

	// Filter by clinic for non-super-admin users
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent template not found"})
	}

	return c.JSON(template)
}

func UpdateConsentTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	// Only super admin and admins can update templates
	if !user.IsSuperAdmin() && !user.HasRole(models.Admin) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	var template models.ConsentTemplate
	query := database.DB.Where("id = ?", id)
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent template not found"})
	}

	var req models.ConsentTemplate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update fields
	template.Name = req.Name
	template.Description = req.Description
	template.Content = req.Content
	template.Category = req.Category
	template.IsActive = req.IsActive
	template.IsDefault = req.IsDefault

	if err := database.DB.Save(&template).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update consent template"})
	}

	return c.JSON(template)
}

func DeleteConsentTemplate(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	// Only super admin and admins can delete templates
	if !user.IsSuperAdmin() && !user.HasRole(models.Admin) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	var template models.ConsentTemplate
	query := database.DB.Where("id = ?", id)
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent template not found"})
	}

	// Soft delete by setting inactive
	template.IsActive = false
	if err := database.DB.Save(&template).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete consent template"})
	}

	return c.JSON(fiber.Map{"message": "Consent template deleted successfully"})
}

// Consent Form Handlers
func GetConsentForms(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	var forms []models.ConsentForm

	query := database.DB.Preload("Patient").Preload("ConsentTemplate").Preload("Doctor")

	// Filter by clinic for non-super-admin users
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	// Filter by patient if provided
	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&forms).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch consent forms"})
	}

	return c.JSON(forms)
}

func CreateConsentForm(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req models.ConsentForm
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.PatientID == 0 || req.ConsentTemplateID == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Patient ID and consent template ID are required"})
	}

	// Set clinic ID for non-super-admin users
	if !user.IsSuperAdmin() {
		req.ClinicID = user.ClinicID
	}

	// Set created by
	req.CreatedByID = user.ID

	// Set status to draft if not provided
	if req.Status == "" {
		req.Status = models.DocStatusDraft
	}

	// Generate content from template with signature areas
	if req.ConsentTemplateID != nil {
		var template models.ConsentTemplate
		if err := database.DB.First(&template, *req.ConsentTemplateID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid consent template"})
		}

		// Load patient for placeholder replacement
		var patient models.Patient
		if err := database.DB.First(&patient, req.PatientID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid patient"})
		}

		// Replace placeholders in template content
		content := template.Content
		content = strings.ReplaceAll(content, "[PATIENT_NAME]", patient.FirstName+" "+patient.LastName)
		content = strings.ReplaceAll(content, "[PATIENT_FIRST_NAME]", patient.FirstName)
		content = strings.ReplaceAll(content, "[PATIENT_LAST_NAME]", patient.LastName)

		currentDate := time.Now().Format("January 2, 2006")
		content = strings.ReplaceAll(content, "[CURRENT_DATE]", currentDate)
		content = strings.ReplaceAll(content, "[TODAY]", currentDate)
		content = strings.ReplaceAll(content, "[DATE]", currentDate)

		if patient.DateOfBirth != nil {
			dob := patient.DateOfBirth.Format("January 2, 2006")
			content = strings.ReplaceAll(content, "[PATIENT_DOB]", dob)
		}

		if patient.Phone != "" {
			content = strings.ReplaceAll(content, "[PATIENT_PHONE]", patient.Phone)
		}

		if patient.Email != "" {
			content = strings.ReplaceAll(content, "[PATIENT_EMAIL]", patient.Email)
		}

		// Handle doctor placeholder if doctor is specified
		if req.DoctorID != nil {
			var doctor models.User
			if err := database.DB.First(&doctor, *req.DoctorID).Error; err == nil {
				doctorName := doctor.FirstName + " " + doctor.LastName
				content = strings.ReplaceAll(content, "[DOCTOR_NAME]", doctorName)
			}
		} else {
			// Use default if no doctor specified
			content = strings.ReplaceAll(content, "[DOCTOR_NAME]", "Doctor")
		}

		// Set the populated content
		req.Content = content
	}

	if err := database.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create consent form"})
	}

	// Load the created form with relations
	if err := database.DB.Preload("Patient").Preload("ConsentTemplate").Preload("Doctor").First(&req, req.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load created consent form"})
	}

	return c.Status(201).JSON(req)
}

func GetConsentForm(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	var form models.ConsentForm
	query := database.DB.Preload("Patient").Preload("ConsentTemplate").Preload("Doctor").Preload("Witness")

	// Filter by clinic for non-super-admin users
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&form, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent form not found"})
	}

	return c.JSON(form)
}

func UpdateConsentForm(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	var form models.ConsentForm
	query := database.DB.Where("id = ?", id)
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&form).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent form not found"})
	}

	var req models.ConsentForm
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update fields
	form.Title = req.Title
	form.Content = req.Content
	form.Status = req.Status
	form.ProcedureDescription = req.ProcedureDescription
	form.Risks = req.Risks
	form.Alternatives = req.Alternatives
	form.PostCareInstructions = req.PostCareInstructions
	form.UnderstandsTreatment = req.UnderstandsTreatment
	form.UnderstandsRisks = req.UnderstandsRisks
	form.ConsentsToTreatment = req.ConsentsToTreatment
	form.HadOpportunityToAsk = req.HadOpportunityToAsk
	form.PatientSignature = req.PatientSignature
	form.WitnessSignature = req.WitnessSignature
	form.WitnessID = req.WitnessID

	// Set timestamps if signatures are provided
	if req.PatientSignature != "" && form.PatientSignedAt == nil {
		now := time.Now()
		form.PatientSignedAt = &now
	}
	if req.WitnessSignature != "" && form.WitnessSignedAt == nil {
		now := time.Now()
		form.WitnessSignedAt = &now
	}

	if err := database.DB.Save(&form).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update consent form"})
	}

	// Load the updated form with relations
	if err := database.DB.Preload("Patient").Preload("ConsentTemplate").Preload("Witness").First(&form, form.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load updated consent form"})
	}

	return c.JSON(form)
}

func SignConsentForm(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := c.Params("id")

	var form models.ConsentForm
	query := database.DB.Where("id = ?", id)
	if !user.IsSuperAdmin() {
		query = query.Where("clinic_id = ?", user.ClinicID)
	}

	if err := query.First(&form).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consent form not found"})
	}

	var req struct {
		PatientSignature string `json:"patient_signature"`
		WitnessSignature string `json:"witness_signature"`
		WitnessID        *uint  `json:"witness_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update signatures
	if req.PatientSignature != "" {
		form.PatientSignature = req.PatientSignature
		now := time.Now()
		form.PatientSignedAt = &now
	}

	if req.WitnessSignature != "" {
		form.WitnessSignature = req.WitnessSignature
		form.WitnessID = req.WitnessID
		now := time.Now()
		form.WitnessSignedAt = &now
	}

	// Update status if both signatures are present
	if form.PatientSignedAt != nil && form.WitnessSignedAt != nil {
		form.Status = models.DocStatusSigned
	} else if form.PatientSignedAt != nil {
		form.Status = models.DocStatusPending
	}

	if err := database.DB.Save(&form).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to sign consent form"})
	}

	return c.JSON(form)
}
