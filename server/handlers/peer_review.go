package handlers

import (
	"net/http"
	"strconv"
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreatePeerReviewCase creates a new peer review case
func CreatePeerReviewCase(c *fiber.Ctx) error {
	// Get user from context
	user := c.Locals("user").(models.User)

	// Only doctors can create peer review cases
	if user.Role != models.Doctor {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Only doctors can create peer review cases",
		})
	}

	var req struct {
		CaseNumber         string `json:"case_number"` // Unique case identifier
		Title              string `json:"title"`
		Description        string `json:"description"`
		Visibility         string `json:"visibility"`           // public, in_clinic, invite_only
		PatientID          *uint  `json:"patient_id,omitempty"` // Optional: link to existing patient
		PatientAge         int    `json:"patient_age"`
		PatientGender      string `json:"patient_gender"`
		PatientBloodType   string `json:"patient_blood_type"`
		ChiefComplaint     string `json:"chief_complaint"`
		MedicalHistory     string `json:"medical_history"`
		DentalHistory      string `json:"dental_history"`
		CurrentMedications string `json:"current_medications"`
		Allergies          string `json:"allergies"`
		DentalChartData    string `json:"dental_chart_data"`
		ShareWith          []uint `json:"share_with,omitempty"` // User IDs to invite (for invite_only)
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.CaseNumber == "" || req.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Case number and title are required",
		})
	}

	// Validate visibility
	visibility := models.VisibilityInviteOnly
	switch req.Visibility {
	case "public":
		visibility = models.VisibilityPublic
	case "in_clinic":
		visibility = models.VisibilityInClinic
	case "invite_only":
		visibility = models.VisibilityInviteOnly
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid visibility level",
		})
	}

	// If patient_id is provided, verify access and get anonymized data
	var originalPatientID *uint
	if req.PatientID != nil {
		var patient models.Patient
		if err := database.DB.Preload("Clinic").First(&patient, *req.PatientID).Error; err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Patient not found",
			})
		}

		// Verify user has access to this patient (same clinic)
		if user.ClinicID == nil || *user.ClinicID != patient.ClinicID {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied to this patient",
			})
		}

		originalPatientID = req.PatientID

		// If dental chart data not provided, get from patient
		if req.DentalChartData == "" {
			var dentalRecords []models.DentalRecord
			database.DB.Where("patient_id = ? AND is_active = ?", patient.ID, true).Find(&dentalRecords)
			anonymizedData := createAnonymizedPatientData(patient, dentalRecords)
			req.DentalChartData = anonymizedData.DentalChart
		}
	}

	// Create peer review case
	peerReviewCase := models.PeerReviewCase{
		CaseNumber:         req.CaseNumber,
		Title:              req.Title,
		Description:        req.Description,
		Visibility:         visibility,
		PatientAge:         req.PatientAge,
		PatientGender:      req.PatientGender,
		PatientBloodType:   req.PatientBloodType,
		ChiefComplaint:     req.ChiefComplaint,
		MedicalHistory:     req.MedicalHistory,
		DentalHistory:      req.DentalHistory,
		CurrentMedications: req.CurrentMedications,
		Allergies:          req.Allergies,
		DentalChartData:    req.DentalChartData,
		Status:             models.PeerReviewStatusOpen,
		CreatedByID:        user.ID,
		ClinicID:           *user.ClinicID,
		OriginalPatientID:  originalPatientID,
	}

	if err := database.DB.Create(&peerReviewCase).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create peer review case",
		})
	}

	// Handle participants based on visibility
	now := time.Now()
	if visibility == models.VisibilityInviteOnly && len(req.ShareWith) > 0 {
		// Add invited participants
		for _, participantID := range req.ShareWith {
			participant := models.PeerReviewParticipant{
				CaseID:      peerReviewCase.ID,
				UserID:      participantID,
				Permission:  models.PermissionComment,
				InvitedByID: user.ID,
				InvitedAt:   now,
				AcceptedAt:  &now, // Auto-accept for now
			}
			database.DB.Create(&participant)
		}
	}

	// Add creator as participant with edit permission
	creatorParticipant := models.PeerReviewParticipant{
		CaseID:      peerReviewCase.ID,
		UserID:      user.ID,
		Permission:  models.PermissionEdit,
		InvitedByID: user.ID,
		InvitedAt:   now,
		AcceptedAt:  &now,
	}
	database.DB.Create(&creatorParticipant)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Peer review case created successfully",
		"case":    peerReviewCase,
	})
}

// GetPeerReviewCases gets all peer review cases the user has access to
func GetPeerReviewCases(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Get query parameters
	status := c.Query("status")
	visibility := c.Query("visibility")
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	var cases []models.PeerReviewCase

	// Combine all accessible cases
	query := database.DB.Raw(`
		SELECT DISTINCT prc.* FROM peer_review_cases prc
		LEFT JOIN peer_review_participants prp ON prc.id = prp.case_id AND prp.user_id = ?
		WHERE (
			prc.visibility = 'public'
			OR (prc.visibility = 'in_clinic' AND prc.clinic_id = ?)
			OR (prc.visibility = 'invite_only' AND prp.user_id IS NOT NULL)
		)
	`, user.ID, *user.ClinicID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}

	var total int64
	query.Model(&models.PeerReviewCase{}).Count(&total)

	if err := query.Preload("CreatedBy").Preload("Clinic").
		Offset(offset).Limit(limit).Find(&cases).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch peer review cases",
		})
	}

	return c.JSON(fiber.Map{
		"cases": cases,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPeerReviewCase gets a specific peer review case
func GetPeerReviewCase(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	caseID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid case ID",
		})
	}

	// Check if user has access to this case
	var participant models.PeerReviewParticipant
	if err := database.DB.Where("case_id = ? AND user_id = ?", caseID, user.ID).First(&participant).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this case",
		})
	}

	var peerCase models.PeerReviewCase
	if err := database.DB.Preload("CreatedBy").Preload("Clinic").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Order("created_at ASC")
		}).
		Preload("Participants", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		First(&peerCase, caseID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Peer review case not found",
		})
	}

	return c.JSON(peerCase)
}

// AddPeerReviewComment adds a comment to a peer review case
func AddPeerReviewComment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	caseID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid case ID",
		})
	}

	var req struct {
		Content     string `json:"content"`
		CommentType string `json:"comment_type"`
		ParentID    *uint  `json:"parent_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Content == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Comment content is required",
		})
	}

	// Check if user has access and comment permission
	var participant models.PeerReviewParticipant
	if err := database.DB.Where("case_id = ? AND user_id = ?", caseID, user.ID).First(&participant).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this case",
		})
	}

	if participant.Permission == models.PermissionView {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to comment on this case",
		})
	}

	// Validate comment type
	commentType := models.CommentTypeGeneral
	switch req.CommentType {
	case "recommendation":
		commentType = models.CommentTypeRecommendation
	case "question":
		commentType = models.CommentTypeQuestion
	case "diagnosis":
		commentType = models.CommentTypeDiagnosis
	}

	comment := models.PeerReviewComment{
		CaseID:      uint(caseID),
		Content:     req.Content,
		CommentType: commentType,
		AuthorID:    user.ID,
		ParentID:    req.ParentID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add comment",
		})
	}

	// Load the created comment with author info
	database.DB.Preload("Author").First(&comment, comment.ID)

	return c.Status(http.StatusCreated).JSON(comment)
}

// UpdatePeerReviewCaseStatus updates the status of a peer review case
func UpdatePeerReviewCaseStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	caseID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid case ID",
		})
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user has edit permission
	var participant models.PeerReviewParticipant
	if err := database.DB.Where("case_id = ? AND user_id = ? AND permission = ?",
		caseID, user.ID, models.PermissionEdit).First(&participant).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this case",
		})
	}

	// Validate status
	var status models.PeerReviewStatus
	switch req.Status {
	case "open":
		status = models.PeerReviewStatusOpen
	case "closed":
		status = models.PeerReviewStatusClosed
	case "resolved":
		status = models.PeerReviewStatusResolved
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	if err := database.DB.Model(&models.PeerReviewCase{}).Where("id = ?", caseID).
		Update("status", status).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update case status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Case status updated successfully",
		"status":  status,
	})
}

// Helper functions

func createAnonymizedPatientData(patient models.Patient, dentalRecords []models.DentalRecord) models.AnonymizedPatientData {
	// Calculate age
	age := 0
	if patient.DateOfBirth != nil {
		age = calculateAge(*patient.DateOfBirth)
	}

	// Get dental chart data (simplified for sharing)
	dentalChart := ""
	if len(dentalRecords) > 0 {
		// Use the most recent dental record
		dentalChart = dentalRecords[0].TeethData
	}

	return models.AnonymizedPatientData{
		Age:            age,
		Gender:         patient.Gender,
		BloodType:      string(patient.BloodType),
		ChiefComplaint: "", // Would need to be filled from appointment data
		MedicalHistory: patient.MedicalConditions,
		DentalHistory:  "", // Would need to be filled from appointment history
		Medications:    patient.CurrentMedications,
		Allergies:      patient.Allergies,
		DentalChart:    dentalChart,
	}
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
