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

	// Original patient data (anonymized)
	PatientAge       int    `json:"patient_age"`
	PatientGender    string `json:"patient_gender" gorm:"size:20"`
	PatientBloodType string `json:"patient_blood_type" gorm:"size:10"`

	// Case details
	ChiefComplaint     string `json:"chief_complaint" gorm:"type:text"`
	MedicalHistory     string `json:"medical_history" gorm:"type:text"`
	DentalHistory      string `json:"dental_history" gorm:"type:text"`
	CurrentMedications string `json:"current_medications" gorm:"type:text"`
	Allergies          string `json:"allergies" gorm:"type:text"`

	// Dental chart data (anonymized - no patient identifiers)
	DentalChartData string `json:"dental_chart_data" gorm:"type:text"` // JSON of tooth conditions

	// Case status
	Status PeerReviewStatus `json:"status" gorm:"type:varchar(20);default:'open'"`

	// Creator and clinic info
	CreatedByID uint   `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   User   `json:"created_by" gorm:"foreignKey:CreatedByID"`
	ClinicID    uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic      Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Original patient reference (for internal use only, not exposed in API)
	OriginalPatientID *uint `json:"-" gorm:"index"`

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

// AnonymizedPatientData represents the anonymized patient information for sharing
type AnonymizedPatientData struct {
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	BloodType      string `json:"blood_type"`
	ChiefComplaint string `json:"chief_complaint"`
	MedicalHistory string `json:"medical_history"`
	DentalHistory  string `json:"dental_history"`
	Medications    string `json:"medications"`
	Allergies      string `json:"allergies"`
	DentalChart    string `json:"dental_chart"` // JSON string of tooth data
}
