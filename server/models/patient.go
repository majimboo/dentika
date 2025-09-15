package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BloodType string

const (
	BloodTypeAPositive  BloodType = "A+"
	BloodTypeANegative  BloodType = "A-"
	BloodTypeBPositive  BloodType = "B+"
	BloodTypeBNegative  BloodType = "B-"
	BloodTypeABPositive BloodType = "AB+"
	BloodTypeABNegative BloodType = "AB-"
	BloodTypeOPositive  BloodType = "O+"
	BloodTypeONegative  BloodType = "O-"
)

type Patient struct {
	ID            uint       `json:"id" gorm:"primarykey"`
	PatientNumber string     `json:"patient_number" gorm:"index;size:50"`
	FirstName     string     `json:"first_name" gorm:"size:100;not null"`
	LastName      string     `json:"last_name" gorm:"size:100;not null"`
	Email         string     `json:"email" gorm:"size:100"`
	Phone         string     `json:"phone" gorm:"size:50"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        string     `json:"gender" gorm:"size:20"`
	Address       string     `json:"address" gorm:"size:500"`

	// Emergency Contact
	EmergencyContactName     string `json:"emergency_contact_name" gorm:"size:200"`
	EmergencyContactPhone    string `json:"emergency_contact_phone" gorm:"size:50"`
	EmergencyContactRelation string `json:"emergency_contact_relation" gorm:"size:100"`

	// Medical Information
	BloodType          BloodType `json:"blood_type" gorm:"size:10"`
	Allergies          string    `json:"allergies" gorm:"type:text"`
	MedicalConditions  string    `json:"medical_conditions" gorm:"type:text"`
	CurrentMedications string    `json:"current_medications" gorm:"type:text"`
	Notes              string    `json:"notes" gorm:"type:text"`

	// Insurance
	InsuranceProvider string `json:"insurance_provider" gorm:"size:200"`
	InsuranceNumber   string `json:"insurance_number" gorm:"size:100"`

	// Status
	IsActive          bool   `json:"is_active" gorm:"default:true"`
	PreferredLanguage string `json:"preferred_language" gorm:"size:50;default:'English'"`
	AvatarPath        string `json:"avatar_path" gorm:"size:500"`

	// Foreign Keys
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Relationships
	Appointments     []Appointment          `json:"appointments,omitempty" gorm:"foreignKey:PatientID"`
	DentalRecords    []DentalRecord         `json:"dental_records,omitempty" gorm:"foreignKey:PatientID"`
	PatientDocuments []PatientDocument      `json:"documents,omitempty" gorm:"foreignKey:PatientID"`
	Diagnoses        []PatientDiagnosis     `json:"diagnoses,omitempty" gorm:"foreignKey:PatientID"`
	TreatmentPlans   []PatientTreatmentPlan `json:"treatment_plans,omitempty" gorm:"foreignKey:PatientID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PatientDocument struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Title       string `json:"title" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`
	FilePath    string `json:"file_path" gorm:"size:500;not null"`
	FileType    string `json:"file_type" gorm:"size:50"`
	FileSize    int64  `json:"file_size"`

	// Foreign Keys
	PatientID    uint    `json:"patient_id" gorm:"not null;index"`
	Patient      Patient `json:"patient" gorm:"foreignKey:PatientID"`
	UploadedByID uint    `json:"uploaded_by_id" gorm:"not null"`
	UploadedBy   User    `json:"uploaded_by" gorm:"foreignKey:UploadedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Patient) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *Patient) GetAge() int {
	if p.DateOfBirth == nil {
		return 0
	}
	return int(time.Since(*p.DateOfBirth).Hours() / 24 / 365.25)
}

func (p *Patient) GeneratePatientNumber(clinicID uint) string {
	now := time.Now()
	return fmt.Sprintf("P%d%04d%02d%05d", clinicID, now.Year(), now.Month(), p.ID)
}

// GenerateUniquePatientNumber generates a unique patient number for a clinic
func GenerateUniquePatientNumber(clinicID uint, db *gorm.DB) (string, error) {
	var count int64
	// Count existing patients for this clinic to get the next number
	if err := db.Model(&Patient{}).Where("clinic_id = ?", clinicID).Count(&count).Error; err != nil {
		return "", err
	}

	now := time.Now()
	// Format: {year}{0000}{sequential_number}
	nextNumber := count + 1
	return fmt.Sprintf("%d%04d", now.Year(), nextNumber), nil
}

// PatientSelfScheduleRequest represents a patient's request for self-scheduling an appointment
type PatientSelfScheduleRequest struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	FirstName string `json:"first_name" gorm:"size:100;not null"`
	LastName  string `json:"last_name" gorm:"size:100;not null"`
	Email     string `json:"email" gorm:"size:100;not null"`
	Phone     string `json:"phone" gorm:"size:50;not null"`

	// Appointment preferences
	PreferredDate   *time.Time `json:"preferred_date"`
	PreferredTime   string     `json:"preferred_time" gorm:"size:10"` // HH:MM format
	Symptoms        string     `json:"symptoms" gorm:"type:text"`     // What's bothering them
	AdditionalNotes string     `json:"additional_notes" gorm:"type:text"`

	// Clinic/Branch information
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`
	BranchID *uint  `json:"branch_id" gorm:"index"`
	Branch   Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`

	// Status tracking
	Status       PatientSelfScheduleStatus `json:"status" gorm:"default:'pending'"`
	ReviewedByID *uint                     `json:"reviewed_by_id"`
	ReviewedBy   User                      `json:"reviewed_by,omitempty" gorm:"foreignKey:ReviewedByID"`
	ReviewNotes  string                    `json:"review_notes" gorm:"type:text"`
	ReviewedAt   *time.Time                `json:"reviewed_at"`

	// Conversion tracking (when converted to actual appointment)
	ConvertedToAppointmentID *uint       `json:"converted_to_appointment_id"`
	ConvertedToAppointment   Appointment `json:"converted_to_appointment,omitempty" gorm:"foreignKey:ConvertedToAppointmentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PatientSelfScheduleStatus string

const (
	PatientScheduleStatusPending   PatientSelfScheduleStatus = "pending"
	PatientScheduleStatusReviewed  PatientSelfScheduleStatus = "reviewed"
	PatientScheduleStatusApproved  PatientSelfScheduleStatus = "approved"
	PatientScheduleStatusRejected  PatientSelfScheduleStatus = "rejected"
	PatientScheduleStatusConverted PatientSelfScheduleStatus = "converted"
)

func (p *PatientSelfScheduleRequest) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *PatientSelfScheduleRequest) GetPreferredDateTime() *time.Time {
	if p.PreferredDate == nil {
		return nil
	}

	// Parse preferred time if available
	if p.PreferredTime != "" {
		dateTimeStr := p.PreferredDate.Format("2006-01-02") + " " + p.PreferredTime
		if parsedTime, err := time.Parse("2006-01-02 15:04", dateTimeStr); err == nil {
			return &parsedTime
		}
	}

	return p.PreferredDate
}
