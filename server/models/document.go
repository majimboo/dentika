package models

import (
	"time"

	"gorm.io/gorm"
)

type DocumentType string
type DocumentStatus string

const (
	DocTypeConsent      DocumentType = "consent"
	DocTypePrescription DocumentType = "prescription"
	DocTypeReport       DocumentType = "report"
	DocTypeInvoice      DocumentType = "invoice"
)

const (
	DocStatusDraft     DocumentStatus = "draft"
	DocStatusPending   DocumentStatus = "pending"
	DocStatusSigned    DocumentStatus = "signed"
	DocStatusCompleted DocumentStatus = "completed"
	DocStatusCancelled DocumentStatus = "cancelled"
)

type ConsentTemplate struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Code        string `json:"code" gorm:"uniqueIndex;size:20"`
	Name        string `json:"name" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`

	// Template content
	Content string `json:"content" gorm:"type:text;not null"`

	// Template metadata
	Category  string `json:"category" gorm:"size:100"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	IsDefault bool   `json:"is_default" gorm:"default:false"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ConsentForm struct {
	ID      uint           `json:"id" gorm:"primarykey"`
	Title   string         `json:"title" gorm:"size:200;not null"`
	Content string         `json:"content" gorm:"type:text;not null"`
	Status  DocumentStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`

	// Template reference
	ConsentTemplateID *uint            `json:"consent_template_id" gorm:"index"`
	ConsentTemplate   *ConsentTemplate `json:"consent_template,omitempty" gorm:"foreignKey:ConsentTemplateID"`

	// Patient and appointment details
	PatientID     uint         `json:"patient_id" gorm:"not null;index"`
	Patient       Patient      `json:"patient" gorm:"foreignKey:PatientID"`
	AppointmentID *uint        `json:"appointment_id" gorm:"index"`
	Appointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`

	// Treating doctor
	DoctorID *uint `json:"doctor_id" gorm:"index"`
	Doctor   *User `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`

	// Procedure details
	ProcedureDescription string `json:"procedure_description" gorm:"type:text"`
	Risks                string `json:"risks" gorm:"type:text"`
	Alternatives         string `json:"alternatives" gorm:"type:text"`
	PostCareInstructions string `json:"post_care_instructions" gorm:"type:text"`

	// Patient agreement checkboxes
	UnderstandsTreatment bool `json:"understands_treatment" gorm:"default:false"`
	UnderstandsRisks     bool `json:"understands_risks" gorm:"default:false"`
	ConsentsToTreatment  bool `json:"consents_to_treatment" gorm:"default:false"`
	HadOpportunityToAsk  bool `json:"had_opportunity_to_ask" gorm:"default:false"`

	// Digital signature
	PatientSignature string     `json:"patient_signature" gorm:"type:text"` // base64 encoded signature
	PatientSignedAt  *time.Time `json:"patient_signed_at"`
	WitnessSignature string     `json:"witness_signature" gorm:"type:text"` // base64 encoded signature
	WitnessSignedAt  *time.Time `json:"witness_signed_at"`
	WitnessID        *uint      `json:"witness_id" gorm:"index"`
	Witness          *User      `json:"witness,omitempty" gorm:"foreignKey:WitnessID"`

	// Generated PDF
	PDFPath string `json:"pdf_path" gorm:"size:500"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Created by
	CreatedByID uint `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   User `json:"created_by" gorm:"foreignKey:CreatedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ConsentTemplate methods
func (ct *ConsentTemplate) GetFormattedName() string {
	if ct.Code != "" {
		return ct.Code + " - " + ct.Name
	}
	return ct.Name
}

// ConsentForm methods
func (cf *ConsentForm) IsFullySigned() bool {
	return cf.PatientSignedAt != nil && cf.WitnessSignedAt != nil
}

func (cf *ConsentForm) IsPatientAgreementComplete() bool {
	return cf.UnderstandsTreatment && cf.UnderstandsRisks &&
		cf.ConsentsToTreatment && cf.HadOpportunityToAsk
}

func (cf *ConsentForm) GeneratePDF() error {
	// TODO: Implement PDF generation logic
	return nil
}
