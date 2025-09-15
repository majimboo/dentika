package models

import (
	"time"

	"gorm.io/gorm"
)

// PeerReviewCase represents a shared patient case for peer review
type PeerReviewCase struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	CaseNumber  string `json:"case_number" gorm:"uniqueIndex;size:50;not null"` // Unique case identifier
	Title       string `json:"title" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`

	// Case visibility
	Visibility PeerReviewVisibility `json:"visibility" gorm:"type:varchar(20);default:'invite_only'"`

	// Case status
	Status PeerReviewStatus `json:"status" gorm:"type:varchar(20);default:'open'"`

	// Creator and clinic info
	CreatedByID uint   `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   User   `json:"created_by" gorm:"foreignKey:CreatedByID"`
	ClinicID    uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic      Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Reference to the original patient and appointment (for internal use only, not exposed in API)
	OriginalPatientID     *uint        `json:"-" gorm:"index"`
	OriginalPatient       *Patient     `json:"-" gorm:"foreignKey:OriginalPatientID"`
	OriginalAppointmentID *uint        `json:"appointment_id" gorm:"index"` // Expose this for frontend use
	OriginalAppointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:OriginalAppointmentID"`

	// Relationships
	Comments     []PeerReviewComment     `json:"comments,omitempty" gorm:"foreignKey:CaseID"`
	Participants []PeerReviewParticipant `json:"participants,omitempty" gorm:"foreignKey:CaseID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PeerReviewStatus defines the status of a peer review case
type PeerReviewStatus string

const (
	PeerReviewStatusOpen     PeerReviewStatus = "open"
	PeerReviewStatusClosed   PeerReviewStatus = "closed"
	PeerReviewStatusResolved PeerReviewStatus = "resolved"
)

// PeerReviewVisibility defines who can see the peer review case
type PeerReviewVisibility string

const (
	VisibilityPublic     PeerReviewVisibility = "public"      // All clinics can see
	VisibilityInClinic   PeerReviewVisibility = "in_clinic"   // Only doctors in same clinic
	VisibilityInviteOnly PeerReviewVisibility = "invite_only" // Only invited doctors
)

// PeerReviewComment represents a comment or recommendation on a peer review case
type PeerReviewComment struct {
	ID     uint           `json:"id" gorm:"primarykey"`
	CaseID uint           `json:"case_id" gorm:"not null;index"`
	Case   PeerReviewCase `json:"case" gorm:"foreignKey:CaseID"`

	// Comment content
	Content     string                `json:"content" gorm:"type:text;not null"`
	CommentType PeerReviewCommentType `json:"comment_type" gorm:"type:varchar(20);default:'comment'"`

	// Author info
	AuthorID uint `json:"author_id" gorm:"not null;index"`
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`

	// Parent comment for threaded discussions
	ParentID *uint              `json:"parent_id,omitempty" gorm:"index"`
	Parent   *PeerReviewComment `json:"parent,omitempty" gorm:"foreignKey:ParentID"`

	// Threaded comments
	Replies []PeerReviewComment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PeerReviewCommentType defines the type of comment
type PeerReviewCommentType string

const (
	CommentTypeGeneral        PeerReviewCommentType = "comment"
	CommentTypeRecommendation PeerReviewCommentType = "recommendation"
	CommentTypeQuestion       PeerReviewCommentType = "question"
	CommentTypeDiagnosis      PeerReviewCommentType = "diagnosis"
)

// PeerReviewParticipant tracks who can access a peer review case
type PeerReviewParticipant struct {
	ID     uint           `json:"id" gorm:"primarykey"`
	CaseID uint           `json:"case_id" gorm:"not null;index"`
	Case   PeerReviewCase `json:"case" gorm:"foreignKey:CaseID"`

	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// Permission level
	Permission PeerReviewPermission `json:"permission" gorm:"type:varchar(20);default:'view'"`

	// Invitation details
	InvitedByID uint `json:"invited_by_id" gorm:"not null"`
	InvitedBy   User `json:"invited_by" gorm:"foreignKey:InvitedByID"`

	InvitedAt  time.Time  `json:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PeerReviewPermission defines access levels for participants
type PeerReviewPermission string

const (
	PermissionView    PeerReviewPermission = "view"
	PermissionComment PeerReviewPermission = "comment"
	PermissionEdit    PeerReviewPermission = "edit"
)

// GetAnonymizedPatientData retrieves anonymized patient and appointment data for sharing
func (prc *PeerReviewCase) GetAnonymizedPatientData() map[string]interface{} {
	if prc.OriginalPatient == nil || prc.OriginalAppointment == nil {
		return map[string]interface{}{}
	}

	patient := prc.OriginalPatient
	appointment := prc.OriginalAppointment

	// Calculate age
	age := 0
	if patient.DateOfBirth != nil {
		now := time.Now()
		birthYear := patient.DateOfBirth.Year()
		age = now.Year() - birthYear
		if now.YearDay() < patient.DateOfBirth.YearDay() {
			age--
		}
	}

	return map[string]interface{}{
		"patient_age":         age,
		"patient_gender":      patient.Gender,
		"patient_blood_type":  patient.BloodType,
		"medical_history":     patient.MedicalConditions,
		"current_medications": patient.CurrentMedications,
		"allergies":           patient.Allergies,
		"chief_complaint":     appointment.Description, // Use Description as chief complaint
		"appointment_date":    appointment.StartTime,   // Use StartTime as appointment date
		"procedures":          appointment.Procedures,
		"diagnoses":           appointment.Diagnoses,
		"pre_notes":           appointment.PreAppointmentNotes,
		"post_notes":          appointment.PostAppointmentNotes,
		"title":               appointment.Title,
	}
}
