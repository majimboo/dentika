package handlers

import (
	"fmt"
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
		CaseNumber    string `json:"case_number"` // Unique case identifier
		Title         string `json:"title"`
		Description   string `json:"description"`
		Visibility    string `json:"visibility"`     // public, in_clinic, invite_only
		PatientID     *uint  `json:"patient_id"`     // Required: patient ID
		AppointmentID *uint  `json:"appointment_id"` // Required: appointment ID
		ShareWith     []uint `json:"share_with,omitempty"` // User IDs to invite (for invite_only)
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.CaseNumber == "" || req.Title == "" || req.PatientID == nil || req.AppointmentID == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Case number, title, patient ID, and appointment ID are required",
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

	// Validate patient and appointment access
	var patient models.Patient
	if err := database.DB.Preload("Clinic").First(&patient, *req.PatientID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	// Verify user has access to this patient (same clinic)
	if user.ClinicID != patient.ClinicID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this patient",
		})
	}

	// Validate appointment exists and belongs to the patient
	var appointment models.Appointment
	if err := database.DB.Preload("Patient").Preload("Procedures").Preload("Diagnoses").
		First(&appointment, *req.AppointmentID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Appointment not found",
		})
	}

	// Verify appointment belongs to the selected patient
	if appointment.PatientID != *req.PatientID {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Appointment does not belong to the selected patient",
		})
	}

	// Verify appointment belongs to user's clinic
	if appointment.Patient.ClinicID != user.ClinicID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this appointment",
		})
	}

	// Create peer review case
	peerReviewCase := models.PeerReviewCase{
		CaseNumber:            req.CaseNumber,
		Title:                 req.Title,
		Description:           req.Description,
		Visibility:            visibility,
		Status:                models.PeerReviewStatusOpen,
		CreatedByID:           user.ID,
		ClinicID:              user.ClinicID,
		OriginalPatientID:     req.PatientID,
		OriginalAppointmentID: req.AppointmentID,
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

	// Build query using GORM instead of raw SQL
	query := database.DB.Model(&models.PeerReviewCase{})

	// Access control: combine all accessible cases
	query = query.Where(
		database.DB.Where("visibility = ?", "public").
			Or(database.DB.Where("visibility = ? AND clinic_id = ?", "in_clinic", user.ClinicID)).
			Or(database.DB.Where("visibility = ? AND id IN (?)", "invite_only",
				database.DB.Table("peer_review_participants").
					Select("case_id").
					Where("user_id = ?", user.ID),
			)),
	)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}

	var total int64
	query.Count(&total)

	if err := query.Preload("CreatedBy").Preload("Clinic").Preload("Comments").
		Offset(offset).Limit(limit).Find(&cases).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch peer review cases",
		})
	}

	// Anonymize case creators for non-case-makers
	for i := range cases {
		if cases[i].CreatedByID != user.ID {
			cases[i].CreatedBy.FirstName = "Case"
			cases[i].CreatedBy.LastName = "Creator"
			cases[i].CreatedBy.Email = ""
		}
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

	// First, load the case to check its visibility
	var peerCase models.PeerReviewCase
	if err := database.DB.First(&peerCase, caseID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Peer review case not found",
		})
	}

	// Check if user has access based on visibility rules
	hasAccess := false

	switch peerCase.Visibility {
	case models.VisibilityPublic:
		// Public cases are accessible to all users
		hasAccess = true
	case models.VisibilityInClinic:
		// In-clinic cases are accessible to users in the same clinic
		hasAccess = (user.ClinicID == peerCase.ClinicID)
	case models.VisibilityInviteOnly:
		// Invite-only cases require explicit participation
		var participant models.PeerReviewParticipant
		err := database.DB.Where("case_id = ? AND user_id = ?", caseID, user.ID).First(&participant).Error
		hasAccess = (err == nil)
	}

	if !hasAccess {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this case",
		})
	}

	// Now load the full case with all relationships
	if err := database.DB.Preload("CreatedBy").Preload("Clinic").
		Preload("OriginalPatient").
		Preload("OriginalAppointment.Procedures").
		Preload("OriginalAppointment.Diagnoses").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").Preload("Author").Preload("Replies", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Author").Order("created_at ASC")
			}).Order("created_at ASC")
		}).
		Preload("Participants", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Where("id = ?", caseID).First(&peerCase).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Peer review case not found",
		})
	}

	// Check if current user is the case maker (creator)
	isCaseMaker := peerCase.CreatedByID == user.ID

	// Anonymize data if user is not the case maker
	if !isCaseMaker {
		// Hide original patient name information
		peerCase.CreatedBy.FirstName = "Case"
		peerCase.CreatedBy.LastName = "Creator"
		peerCase.CreatedBy.Email = ""

		// Anonymize comments authors (except case maker's own comments)
		for i := range peerCase.Comments {
			if peerCase.Comments[i].AuthorID != peerCase.CreatedByID && peerCase.Comments[i].AuthorID != user.ID {
				peerCase.Comments[i].Author.FirstName = "Dr."
				peerCase.Comments[i].Author.LastName = getAnonymousName(peerCase.Comments[i].AuthorID)
				peerCase.Comments[i].Author.Email = ""
			} else if peerCase.Comments[i].AuthorID == peerCase.CreatedByID {
				// Case maker's comments show as "Case Creator"
				peerCase.Comments[i].Author.FirstName = "Case"
				peerCase.Comments[i].Author.LastName = "Creator"
				peerCase.Comments[i].Author.Email = ""
			}

			// Also anonymize replies
			for j := range peerCase.Comments[i].Replies {
				if peerCase.Comments[i].Replies[j].AuthorID != peerCase.CreatedByID && peerCase.Comments[i].Replies[j].AuthorID != user.ID {
					peerCase.Comments[i].Replies[j].Author.FirstName = "Dr."
					peerCase.Comments[i].Replies[j].Author.LastName = getAnonymousName(peerCase.Comments[i].Replies[j].AuthorID)
					peerCase.Comments[i].Replies[j].Author.Email = ""
				} else if peerCase.Comments[i].Replies[j].AuthorID == peerCase.CreatedByID {
					// Case maker's replies show as "Case Creator"
					peerCase.Comments[i].Replies[j].Author.FirstName = "Case"
					peerCase.Comments[i].Replies[j].Author.LastName = "Creator"
					peerCase.Comments[i].Replies[j].Author.Email = ""
				}
			}
		}

		// Anonymize participants (except case maker and current user)
		for i := range peerCase.Participants {
			if peerCase.Participants[i].UserID != peerCase.CreatedByID && peerCase.Participants[i].UserID != user.ID {
				peerCase.Participants[i].User.FirstName = "Dr."
				peerCase.Participants[i].User.LastName = getAnonymousName(peerCase.Participants[i].UserID)
				peerCase.Participants[i].User.Email = ""
			} else if peerCase.Participants[i].UserID == peerCase.CreatedByID {
				// Case maker shows as "Case Creator"
				peerCase.Participants[i].User.FirstName = "Case"
				peerCase.Participants[i].User.LastName = "Creator"
				peerCase.Participants[i].User.Email = ""
			}
		}
	}

	// Include anonymized patient data from appointment
	response := fiber.Map{
		"case":          peerCase,
		"is_case_maker": isCaseMaker,
	}

	// Add anonymized patient/appointment data if available
	if peerCase.OriginalPatient != nil && peerCase.OriginalAppointment != nil {
		response["patient_data"] = peerCase.GetAnonymizedPatientData()
	}

	return c.JSON(response)
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

	// First, load the case to check visibility and access
	var peerCase models.PeerReviewCase
	if err := database.DB.First(&peerCase, caseID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Peer review case not found",
		})
	}

	// Check if user has access based on visibility rules
	hasAccess := false

	switch peerCase.Visibility {
	case models.VisibilityPublic:
		// Public cases are accessible to all users
		hasAccess = true
	case models.VisibilityInClinic:
		// In-clinic cases are accessible to users in the same clinic
		hasAccess = (user.ClinicID == peerCase.ClinicID)
	case models.VisibilityInviteOnly:
		// Invite-only cases require explicit participation
		var participant models.PeerReviewParticipant
		err := database.DB.Where("case_id = ? AND user_id = ?", caseID, user.ID).First(&participant).Error
		if err == nil {
			hasAccess = true
			// For invite-only cases, check if user has comment permission
			if participant.Permission == models.PermissionView {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{
					"error": "You don't have permission to comment on this case",
				})
			}
		}
	}

	if !hasAccess {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this case",
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

// getAnonymousName generates a consistent anonymous name based on user ID
func getAnonymousName(userID uint) string {
	// Generate anonymous names like "Anonymous A", "Anonymous B", etc.
	// This provides consistent anonymization per user within a case
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	index := (userID % uint(len(letters)))
	return fmt.Sprintf("Anonymous %c", letters[index])
}
