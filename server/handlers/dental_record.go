package handlers

import (
	"strconv"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type UpdateToothRequest struct {
	ToothNumber       string                    `json:"tooth_number"`
	Condition         models.ToothCondition     `json:"condition"`
	Surfaces          []string                  `json:"surfaces"`
	SurfaceConditions []models.SurfaceCondition `json:"surface_conditions"`
	Notes             string                    `json:"notes"`
	Reason            string                    `json:"reason"` // Reason for the change
}

func GetPatientDentalRecords(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	patientID, err := strconv.ParseUint(c.Params("patient_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid patient ID"})
	}

	// Verify patient belongs to accessible clinic
	var patient models.Patient
	if err := database.DB.First(&patient, patientID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Patient not found"})
	}

	if !user.IsSuperAdmin() && user.ClinicID != patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var dentalRecords []models.DentalRecord
	if err := database.DB.Where("patient_id = ?", patientID).Find(&dentalRecords).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dental records"})
	}

	// Check for missing dental records and create them
	permanentExists := false
	primaryExists := false
	for _, record := range dentalRecords {
		if record.RecordType == models.ToothTypePermanent {
			permanentExists = true
		}
		if record.RecordType == models.ToothTypePrimary {
			primaryExists = true
		}
	}

	// Create missing permanent teeth record
	if !permanentExists {
		permanentRecord := models.DentalRecord{
			PatientID:  uint(patientID),
			ClinicID:   patient.ClinicID,
			RecordType: models.ToothTypePermanent,
			IsActive:   true,
		}

		// Initialize with complete FDI tooth data
		_, err := permanentRecord.GetTeethData()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to initialize permanent teeth data"})
		}

		if err := database.DB.Create(&permanentRecord).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create permanent teeth record"})
		}

		dentalRecords = append(dentalRecords, permanentRecord)
	}

	// Create missing primary teeth record
	if !primaryExists {
		primaryRecord := models.DentalRecord{
			PatientID:  uint(patientID),
			ClinicID:   patient.ClinicID,
			RecordType: models.ToothTypePrimary,
			IsActive:   false, // Primary teeth start as inactive
		}

		// Initialize with complete primary tooth data
		_, err := primaryRecord.GetTeethData()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to initialize primary teeth data"})
		}

		if err := database.DB.Create(&primaryRecord).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create primary teeth record"})
		}

		dentalRecords = append(dentalRecords, primaryRecord)
	}

	// Parse teeth data for each record and check for completeness
	var recordsWithTeeth []fiber.Map
	for i := range dentalRecords {
		record := &dentalRecords[i] // Use pointer to modify if needed

		teethData, err := record.GetTeethData()

		// Check if all expected teeth are present
		expectedCount := 32
		if record.RecordType == models.ToothTypePrimary {
			expectedCount = 20
		}

		needsReinit := err != nil || len(teethData) != expectedCount

		if needsReinit {
			// Reinitialize with complete tooth data by clearing TeethData and getting fresh data
			record.TeethData = ""
			_, initErr := record.GetTeethData()
			if initErr != nil {
				continue // Skip this record
			}
			// Save the reinitialized data
			if err := database.DB.Save(record).Error; err != nil {
				continue // Skip this record
			}
		}

		recordsWithTeeth = append(recordsWithTeeth, fiber.Map{
			"id":          record.ID,
			"patient_id":  record.PatientID,
			"record_type": record.RecordType,
			"is_active":   record.IsActive,
			"teeth_data":  teethData,
			"created_at":  record.CreatedAt,
			"updated_at":  record.UpdatedAt,
		})
	}

	return c.JSON(recordsWithTeeth)
}

func GetDentalRecord(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	// Check access
	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	teethData, err := dentalRecord.GetTeethData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse teeth data"})
	}

	return c.JSON(fiber.Map{
		"id":          dentalRecord.ID,
		"patient_id":  dentalRecord.PatientID,
		"patient":     dentalRecord.Patient,
		"record_type": dentalRecord.RecordType,
		"is_active":   dentalRecord.IsActive,
		"teeth_data":  teethData,
		"created_at":  dentalRecord.CreatedAt,
		"updated_at":  dentalRecord.UpdatedAt,
	})
}

func UpdateToothCondition(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	// Check access and permissions
	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can modify dental records
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can modify dental records"})
	}

	var req UpdateToothRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.ToothNumber == "" || req.Condition == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tooth number and condition are required"})
	}

	// Get current teeth data to find previous condition
	teethData, err := dentalRecord.GetTeethData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse current teeth data"})
	}

	var previousCondition models.ToothCondition
	for _, tooth := range teethData {
		if tooth.ToothNumber == req.ToothNumber {
			previousCondition = tooth.Condition
			break
		}
	}

	// Update tooth condition
	if err := dentalRecord.UpdateToothCondition(req.ToothNumber, req.Condition, req.Surfaces, req.SurfaceConditions, req.Notes, user.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update tooth condition"})
	}

	// Save the record
	if err := database.DB.Save(&dentalRecord).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save dental record"})
	}

	// Create history entry
	history := models.DentalRecordHistory{
		DentalRecordID:    dentalRecord.ID,
		ToothNumber:       req.ToothNumber,
		PreviousCondition: previousCondition,
		NewCondition:      req.Condition,
		ChangeReason:      req.Reason,
		ChangedByID:       user.ID,
	}

	database.DB.Create(&history)

	// Return updated teeth data
	updatedTeethData, _ := dentalRecord.GetTeethData()

	return c.JSON(fiber.Map{
		"message":    "Tooth condition updated successfully",
		"teeth_data": updatedTeethData,
	})
}

func GetDentalRecordHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	// Verify access to dental record
	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var history []models.DentalRecordHistory
	if err := database.DB.Preload("ChangedBy").Preload("Appointment").
		Where("dental_record_id = ?", recordID).
		Order("created_at DESC").Find(&history).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch history"})
	}

	return c.JSON(history)
}

func GetToothHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	toothNumber := c.Query("tooth_number")
	if toothNumber == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tooth number is required"})
	}

	// Verify access to dental record
	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var history []models.DentalRecordHistory
	if err := database.DB.Preload("ChangedBy").Preload("Appointment").
		Where("dental_record_id = ? AND tooth_number = ?", recordID, toothNumber).
		Order("created_at DESC").Find(&history).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tooth history"})
	}

	return c.JSON(history)
}

func ActivateDentalRecord(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	// Check access and permissions
	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can activate records
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can activate dental records"})
	}

	// Deactivate other records of same type for this patient
	database.DB.Model(&models.DentalRecord{}).
		Where("patient_id = ? AND record_type = ? AND id != ?",
			dentalRecord.PatientID, dentalRecord.RecordType, dentalRecord.ID).
		Update("is_active", false)

	// Activate this record
	dentalRecord.IsActive = true
	if err := database.DB.Save(&dentalRecord).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to activate dental record"})
	}

	return c.JSON(fiber.Map{"message": "Dental record activated successfully"})
}

func BulkUpdateTeeth(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	recordID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var dentalRecord models.DentalRecord
	if err := database.DB.Preload("Patient").First(&dentalRecord, recordID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dental record not found"})
	}

	// Check access and permissions
	if !user.IsSuperAdmin() && user.ClinicID != dentalRecord.Patient.ClinicID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// Only doctors can modify dental records
	if !user.IsSuperAdmin() && !user.HasRole(models.Doctor) {
		return c.Status(403).JSON(fiber.Map{"error": "Only doctors can modify dental records"})
	}

	var req struct {
		Updates []UpdateToothRequest `json:"updates"`
		Reason  string               `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(req.Updates) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No updates provided"})
	}

	// Get current teeth data for history tracking
	currentTeethData, err := dentalRecord.GetTeethData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse current teeth data"})
	}

	// Process each update
	for _, update := range req.Updates {
		if update.ToothNumber == "" || update.Condition == "" {
			continue
		}

		// Find previous condition
		var previousCondition models.ToothCondition
		for _, tooth := range currentTeethData {
			if tooth.ToothNumber == update.ToothNumber {
				previousCondition = tooth.Condition
				break
			}
		}

		// Update tooth condition
		if err := dentalRecord.UpdateToothCondition(update.ToothNumber, update.Condition, update.Surfaces, update.SurfaceConditions, update.Notes, user.ID); err != nil {
			continue // Log error but continue with other updates
		}

		// Create history entry
		reason := update.Reason
		if reason == "" {
			reason = req.Reason
		}

		history := models.DentalRecordHistory{
			DentalRecordID:    dentalRecord.ID,
			ToothNumber:       update.ToothNumber,
			PreviousCondition: previousCondition,
			NewCondition:      update.Condition,
			ChangeReason:      reason,
			ChangedByID:       user.ID,
		}

		database.DB.Create(&history)
	}

	// Save the record
	if err := database.DB.Save(&dentalRecord).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save dental record"})
	}

	// Return updated teeth data
	updatedTeethData, _ := dentalRecord.GetTeethData()

	return c.JSON(fiber.Map{
		"message":    "Bulk update completed successfully",
		"teeth_data": updatedTeethData,
	})
}
