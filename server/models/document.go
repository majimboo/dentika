package models

import (
	"fmt"
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

type ConsentForm struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Title         string         `json:"title" gorm:"size:200;not null"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	Status        DocumentStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
	
	// Patient and appointment details
	PatientID     uint           `json:"patient_id" gorm:"not null;index"`
	Patient       Patient        `json:"patient" gorm:"foreignKey:PatientID"`
	AppointmentID *uint          `json:"appointment_id" gorm:"index"`
	Appointment   *Appointment   `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	
	// Procedure details
	ProcedureDescription string `json:"procedure_description" gorm:"type:text"`
	Risks               string `json:"risks" gorm:"type:text"`
	Alternatives        string `json:"alternatives" gorm:"type:text"`
	PostCareInstructions string `json:"post_care_instructions" gorm:"type:text"`
	
	// Digital signature
	PatientSignature    string     `json:"patient_signature" gorm:"type:text"`      // base64 encoded signature
	PatientSignedAt     *time.Time `json:"patient_signed_at"`
	WitnessSignature    string     `json:"witness_signature" gorm:"type:text"`     // base64 encoded signature
	WitnessSignedAt     *time.Time `json:"witness_signed_at"`
	WitnessID           *uint      `json:"witness_id" gorm:"index"`
	Witness             *User      `json:"witness,omitempty" gorm:"foreignKey:WitnessID"`
	
	// Generated PDF
	PDFPath             string     `json:"pdf_path" gorm:"size:500"`
	
	// Created by
	CreatedByID   uint           `json:"created_by_id" gorm:"not null;index"`
	CreatedBy     User           `json:"created_by" gorm:"foreignKey:CreatedByID"`
	
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Prescription struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	PrescriptionNumber string      `json:"prescription_number" gorm:"uniqueIndex;size:50"`
	Status          DocumentStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
	
	// Patient and appointment details
	PatientID       uint           `json:"patient_id" gorm:"not null;index"`
	Patient         Patient        `json:"patient" gorm:"foreignKey:PatientID"`
	AppointmentID   *uint          `json:"appointment_id" gorm:"index"`
	Appointment     *Appointment   `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	
	// Doctor details
	DoctorID        uint           `json:"doctor_id" gorm:"not null;index"`
	Doctor          User           `json:"doctor" gorm:"foreignKey:DoctorID"`
	DoctorSignature string         `json:"doctor_signature" gorm:"type:text"`       // base64 encoded signature
	SignedAt        *time.Time     `json:"signed_at"`
	
	// Clinic header information
	ClinicID        uint           `json:"clinic_id" gorm:"not null;index"`
	Clinic          Clinic         `json:"clinic" gorm:"foreignKey:ClinicID"`
	
	// Prescription details
	Instructions    string         `json:"instructions" gorm:"type:text"`
	Notes           string         `json:"notes" gorm:"type:text"`
	
	// Relationships
	Medications     []PrescriptionMedication `json:"medications,omitempty" gorm:"foreignKey:PrescriptionID"`
	
	// Generated PDF
	PDFPath         string         `json:"pdf_path" gorm:"size:500"`
	
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type PrescriptionMedication struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	PrescriptionID uint         `json:"prescription_id" gorm:"not null;index"`
	Prescription   Prescription `json:"prescription" gorm:"foreignKey:PrescriptionID"`
	
	MedicationName string       `json:"medication_name" gorm:"size:200;not null"`
	Dosage         string       `json:"dosage" gorm:"size:100;not null"`        // e.g., "500mg"
	Frequency      string       `json:"frequency" gorm:"size:100;not null"`     // e.g., "twice daily"
	Duration       string       `json:"duration" gorm:"size:100"`               // e.g., "7 days"
	Instructions   string       `json:"instructions" gorm:"type:text"`          // e.g., "take with food"
	Quantity       int          `json:"quantity"`                               // number of pills/units
	
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// ConsentForm methods
func (cf *ConsentForm) IsFullySigned() bool {
	return cf.PatientSignedAt != nil && cf.WitnessSignedAt != nil
}

func (cf *ConsentForm) GeneratePDF() error {
	// TODO: Implement PDF generation logic
	return nil
}

// Prescription methods
func (p *Prescription) IsSigned() bool {
	return p.SignedAt != nil && p.DoctorSignature != ""
}

func (p *Prescription) GeneratePrescriptionNumber(clinicID uint) string {
	now := time.Now()
	return fmt.Sprintf("RX%d%04d%02d%05d", clinicID, now.Year(), now.Month(), p.ID)
}

func (p *Prescription) GeneratePDF() error {
	// TODO: Implement PDF generation logic
	return nil
}