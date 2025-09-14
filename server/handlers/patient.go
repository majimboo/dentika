package handlers

import (
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreatePatientRequest struct {
	FirstName                string     `json:"first_name"`
	LastName                 string     `json:"last_name"`
	Email                    string     `json:"email"`
	Phone                    string     `json:"phone"`
	DateOfBirth              *time.Time `json:"date_of_birth"`
	Gender                   string     `json:"gender"`
	Address                  string     `json:"address"`
	EmergencyContactName     string     `json:"emergency_contact_name"`
	EmergencyContactPhone    string     `json:"emergency_contact_phone"`
	EmergencyContactRelation string     `json:"emergency_contact_relation"`
	BloodType                string     `json:"blood_type"`
	Allergies                string     `json:"allergies"`
	MedicalConditions        string     `json:"medical_conditions"`
	CurrentMedications       string     `json:"current_medications"`
	Notes                    string     `json:"notes"`
	InsuranceProvider        string     `json:"insurance_provider"`
	InsuranceNumber          string     `json:"insurance_number"`
	PreferredLanguage        string     `json:"preferred_language"`
	AvatarPath               string     `json:"avatar_path"`
}

func GetPatients(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Get clinic ID from user context or URL parameter
	var clinicID uint
	if user.IsSuperAdmin() {
		// Super admin can optionally filter by clinic
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			}
		}
	} else {
		clinicID = user.ClinicID
	}

	var patients []models.Patient
	query := database.DB.Preload("Clinic")

	if clinicID > 0 {
		query = query.Where("clinic_id = ?", clinicID)
	} else if !user.IsSuperAdmin() {
		// Non-super admin must have clinic restriction
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Add search functionality
	if search := c.Query("search"); search != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Add status filter
	if status := c.Query("status"); status != "" {
		if status == "active" {
			query = query.Where("is_active = ?", true)
		} else if status == "inactive" {
			query = query.Where("is_active = ?", false)
		}
	}

	// Add sorting
	sortBy := c.Query("sort", "name")
	switch sortBy {
	case "name":
		query = query.Order("first_name ASC, last_name ASC")
	case "created_at":
		query = query.Order("created_at DESC")
	case "last_visit":
		// For last_visit, we might need to join with appointments, but for now sort by created_at
		query = query.Order("created_at DESC")
	default:
		query = query.Order("first_name ASC, last_name ASC")
	}

	// Add pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&patients).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch patients"})
	}

	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.Patient{})
	if clinicID > 0 {
		countQuery = countQuery.Where("clinic_id = ?", clinicID)
	}
	if search := c.Query("search"); search != "" {
		countQuery = countQuery.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status := c.Query("status"); status != "" {
		if status == "active" {
			countQuery = countQuery.Where("is_active = ?", true)
		} else if status == "inactive" {
			countQuery = countQuery.Where("is_active = ?", false)
		}
	}
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"patients": patients,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func GetPatient(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var patient models.Patient
	query := database.DB.Preload("Clinic").Preload("Appointments").Preload("DentalRecords")

	if err := query.First(&patient, patientID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	return c.JSON(patient)
}

func CreatePatient(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req CreatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "First name and last name are required"})
	}

	// Determine clinic ID
	var clinicID uint
	if user.IsSuperAdmin() {
		// Super admin can create for any clinic
		if clinicIDParam := c.Query("clinic_id"); clinicIDParam != "" {
			if id, err := strconv.ParseUint(clinicIDParam, 10, 32); err == nil {
				clinicID = uint(id)
			} else {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
			}
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Clinic ID required for super admin"})
		}
	} else {
		clinicID = user.ClinicID
	}

	// Check if patient with same email already exists in clinic
	if req.Email != "" {
		var existingPatient models.Patient
		if err := database.DB.Where("email = ? AND clinic_id = ?", req.Email, clinicID).First(&existingPatient).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Patient with this email already exists"})
		}
	}

	// Generate unique patient number before creating patient
	patientNumber, err := models.GenerateUniquePatientNumber(clinicID, database.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate patient number"})
	}

	patient := models.Patient{
		PatientNumber:            patientNumber,
		FirstName:                req.FirstName,
		LastName:                 req.LastName,
		Email:                    req.Email,
		Phone:                    req.Phone,
		DateOfBirth:              req.DateOfBirth,
		Gender:                   req.Gender,
		Address:                  req.Address,
		EmergencyContactName:     req.EmergencyContactName,
		EmergencyContactPhone:    req.EmergencyContactPhone,
		EmergencyContactRelation: req.EmergencyContactRelation,
		Allergies:                req.Allergies,
		MedicalConditions:        req.MedicalConditions,
		CurrentMedications:       req.CurrentMedications,
		Notes:                    req.Notes,
		InsuranceProvider:        req.InsuranceProvider,
		InsuranceNumber:          req.InsuranceNumber,
		PreferredLanguage:        req.PreferredLanguage,
		AvatarPath:               req.AvatarPath,
		IsActive:                 true,
		ClinicID:                 clinicID,
	}

	// Set blood type if provided
	if req.BloodType != "" {
		patient.BloodType = models.BloodType(req.BloodType)
	}

	if req.PreferredLanguage == "" {
		patient.PreferredLanguage = "English"
	}

	if err := database.DB.Create(&patient).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create patient"})
	}

	// Create initial dental records
	go createInitialDentalRecords(patient.ID, clinicID)

	// Reload with relationships
	database.DB.Preload("Clinic").First(&patient, patient.ID)

	return c.Status(201).JSON(patient)
}

func UpdatePatient(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var patient models.Patient
	if err := database.DB.First(&patient, patientID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var req CreatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update patient fields
	if req.FirstName != "" {
		patient.FirstName = req.FirstName
	}
	if req.LastName != "" {
		patient.LastName = req.LastName
	}
	if req.Email != "" {
		// Check for duplicate email in same clinic
		var existingPatient models.Patient
		if err := database.DB.Where("email = ? AND clinic_id = ? AND id != ?", req.Email, patient.ClinicID, patient.ID).First(&existingPatient).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Patient with this email already exists"})
		}
		patient.Email = req.Email
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		patient.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		patient.Gender = req.Gender
	}
	if req.Address != "" {
		patient.Address = req.Address
	}
	if req.EmergencyContactName != "" {
		patient.EmergencyContactName = req.EmergencyContactName
	}
	if req.EmergencyContactPhone != "" {
		patient.EmergencyContactPhone = req.EmergencyContactPhone
	}
	if req.EmergencyContactRelation != "" {
		patient.EmergencyContactRelation = req.EmergencyContactRelation
	}
	if req.BloodType != "" {
		patient.BloodType = models.BloodType(req.BloodType)
	}
	if req.Allergies != "" {
		patient.Allergies = req.Allergies
	}
	if req.MedicalConditions != "" {
		patient.MedicalConditions = req.MedicalConditions
	}
	if req.CurrentMedications != "" {
		patient.CurrentMedications = req.CurrentMedications
	}
	if req.Notes != "" {
		patient.Notes = req.Notes
	}
	if req.InsuranceProvider != "" {
		patient.InsuranceProvider = req.InsuranceProvider
	}
	if req.InsuranceNumber != "" {
		patient.InsuranceNumber = req.InsuranceNumber
	}
	if req.PreferredLanguage != "" {
		patient.PreferredLanguage = req.PreferredLanguage
	}
	if req.AvatarPath != "" {
		patient.AvatarPath = req.AvatarPath
	}

	if err := database.DB.Save(&patient).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update patient"})
	}

	// Reload with relationships
	database.DB.Preload("Clinic").First(&patient, patient.ID)

	return c.JSON(patient)
}

func DeactivatePatient(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	var patient models.Patient
	if err := database.DB.First(&patient, patientID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
	}

	// Check access - only admin or super admin can deactivate
	if !user.IsSuperAdmin() && !user.HasRole(models.Admin) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	if !user.IsSuperAdmin() && user.ClinicID != patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	patient.IsActive = false
	if err := database.DB.Save(&patient).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to deactivate patient"})
	}

	return c.JSON(fiber.Map{"message": "Patient deactivated successfully"})
}

func createInitialDentalRecords(patientID uint, clinicID uint) {
	// Create permanent teeth record
	permanentRecord := models.DentalRecord{
		PatientID:  patientID,
		ClinicID:   clinicID,
		RecordType: models.ToothTypePermanent,
		IsActive:   true,
	}

	if err := database.DB.Create(&permanentRecord).Error; err != nil {
		return // Log error in production
	}

	// Initialize with default teeth data
	teethData, err := permanentRecord.GetTeethData()
	if err == nil {
		permanentRecord.SetTeethData(teethData)
		database.DB.Save(&permanentRecord)
	}

	// Create primary teeth record for children (can be activated later if needed)
	primaryRecord := models.DentalRecord{
		PatientID:  patientID,
		ClinicID:   clinicID,
		RecordType: models.ToothTypePrimary,
		IsActive:   false, // Inactive by default
	}

	if err := database.DB.Create(&primaryRecord).Error; err != nil {
		return // Log error in production
	}

	teethData, err = primaryRecord.GetTeethData()
	if err == nil {
		primaryRecord.SetTeethData(teethData)
		database.DB.Save(&primaryRecord)
	}
}
