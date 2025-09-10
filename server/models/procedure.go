package models

import (
	"time"

	"gorm.io/gorm"
)

type ProcedureTemplate struct {
	ID                uint    `json:"id" gorm:"primarykey"`
	Code              string  `json:"code" gorm:"uniqueIndex;size:20"`
	Name              string  `json:"name" gorm:"size:200;not null"`
	Description       string  `json:"description" gorm:"type:text"`
	Category          string  `json:"category" gorm:"size:100"`
	EstimatedDuration int     `json:"estimated_duration"` // in minutes
	DefaultCost       float64 `json:"default_cost" gorm:"type:decimal(10,2)"`
	IsActive          bool    `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AppointmentProcedure struct {
	ID                  uint              `json:"id" gorm:"primarykey"`
	ProcedureTemplateID uint              `json:"procedure_template_id" gorm:"not null;index"`
	ProcedureTemplate   ProcedureTemplate `json:"procedure_template" gorm:"foreignKey:ProcedureTemplateID"`

	AppointmentID uint        `json:"appointment_id" gorm:"not null;index"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`

	// Procedure specific details
	ToothNumber string  `json:"tooth_number" gorm:"size:10"`
	Surface     string  `json:"surface" gorm:"size:50"` // mesial, distal, occlusal, etc.
	Notes       string  `json:"notes" gorm:"type:text"`
	Cost        float64 `json:"cost" gorm:"type:decimal(10,2)"`
	Status      string  `json:"status" gorm:"size:50;default:'planned'"` // planned, in_progress, completed, cancelled

	// Timing
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	// Performed by
	PerformedByID uint `json:"performed_by_id" gorm:"not null;index"`
	PerformedBy   User `json:"performed_by" gorm:"foreignKey:PerformedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DiagnosisTemplate struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Code        string `json:"code" gorm:"uniqueIndex;size:20"`
	Name        string `json:"name" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`
	Category    string `json:"category" gorm:"size:100"`
	Severity    string `json:"severity" gorm:"size:50"` // mild, moderate, severe
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AppointmentDiagnosis struct {
	ID                  uint              `json:"id" gorm:"primarykey"`
	DiagnosisTemplateID uint              `json:"diagnosis_template_id" gorm:"not null;index"`
	DiagnosisTemplate   DiagnosisTemplate `json:"diagnosis_template" gorm:"foreignKey:DiagnosisTemplateID"`

	AppointmentID uint        `json:"appointment_id" gorm:"not null;index"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`

	// Diagnosis specific details
	ToothNumber   string `json:"tooth_number" gorm:"size:10"`
	Surface       string `json:"surface" gorm:"size:50"`
	Notes         string `json:"notes" gorm:"type:text"`
	Severity      string `json:"severity" gorm:"size:50"`
	TreatmentPlan string `json:"treatment_plan" gorm:"type:text"`

	// Diagnosed by
	DiagnosedByID uint `json:"diagnosed_by_id" gorm:"not null;index"`
	DiagnosedBy   User `json:"diagnosed_by" gorm:"foreignKey:DiagnosedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (pt *ProcedureTemplate) GetFormattedName() string {
	if pt.Code != "" {
		return pt.Code + " - " + pt.Name
	}
	return pt.Name
}

func (dt *DiagnosisTemplate) GetFormattedName() string {
	if dt.Code != "" {
		return dt.Code + " - " + dt.Name
	}
	return dt.Name
}

func (ap *AppointmentProcedure) GetDuration() *time.Duration {
	if ap.StartTime != nil && ap.EndTime != nil {
		duration := ap.EndTime.Sub(*ap.StartTime)
		return &duration
	}
	return nil
}

// PatientDiagnosis - Patient-specific diagnosis records
type PatientDiagnosis struct {
	ID                  uint              `json:"id" gorm:"primarykey"`
	DiagnosisTemplateID uint              `json:"diagnosis_template_id" gorm:"not null;index"`
	DiagnosisTemplate   DiagnosisTemplate `json:"diagnosis_template" gorm:"foreignKey:DiagnosisTemplateID"`

	PatientID uint    `json:"patient_id" gorm:"not null;index"`
	Patient   Patient `json:"patient" gorm:"foreignKey:PatientID"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Diagnosis specific details
	ToothNumber string `json:"tooth_number" gorm:"size:10"`
	Surface     string `json:"surface" gorm:"size:50"`
	Notes       string `json:"notes" gorm:"type:text"`
	Severity    string `json:"severity" gorm:"size:50"`
	Status      string `json:"status" gorm:"size:50;default:'active'"` // active, resolved, inactive

	// Associated appointment (optional - diagnosis might be made outside appointment)
	AppointmentID *uint        `json:"appointment_id" gorm:"index"`
	Appointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`

	// Diagnosed by
	DiagnosedByID uint `json:"diagnosed_by_id" gorm:"not null;index"`
	DiagnosedBy   User `json:"diagnosed_by" gorm:"foreignKey:DiagnosedByID"`

	// Dates
	DiagnosedAt time.Time  `json:"diagnosed_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PatientTreatmentPlan - Patient-specific treatment plans
type PatientTreatmentPlan struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	PatientID uint    `json:"patient_id" gorm:"not null;index"`
	Patient   Patient `json:"patient" gorm:"foreignKey:PatientID"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Treatment plan details
	Title       string `json:"title" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`
	Priority    string `json:"priority" gorm:"size:50;default:'medium'"` // low, medium, high, urgent
	Status      string `json:"status" gorm:"size:50;default:'active'"`   // active, completed, cancelled, on_hold

	// Associated diagnosis (optional)
	DiagnosisID *uint             `json:"diagnosis_id" gorm:"index"`
	Diagnosis   *PatientDiagnosis `json:"diagnosis,omitempty" gorm:"foreignKey:DiagnosisID"`

	// Estimated details
	EstimatedCost     float64    `json:"estimated_cost" gorm:"type:decimal(10,2)"`
	EstimatedVisits   int        `json:"estimated_visits" gorm:"default:1"`
	EstimatedDuration int        `json:"estimated_duration"` // in days
	StartDate         *time.Time `json:"start_date"`
	TargetCompletion  *time.Time `json:"target_completion"`

	// Progress tracking
	CompletedVisits int     `json:"completed_visits" gorm:"default:0"`
	ActualCost      float64 `json:"actual_cost" gorm:"type:decimal(10,2);default:0"`

	// Created by
	CreatedByID uint `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   User `json:"created_by" gorm:"foreignKey:CreatedByID"`

	// Dates
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TreatmentPlanProcedure - Links procedures to treatment plans
type TreatmentPlanProcedure struct {
	ID uint `json:"id" gorm:"primarykey"`

	TreatmentPlanID uint                 `json:"treatment_plan_id" gorm:"not null;index"`
	TreatmentPlan   PatientTreatmentPlan `json:"treatment_plan" gorm:"foreignKey:TreatmentPlanID"`

	ProcedureTemplateID uint              `json:"procedure_template_id" gorm:"not null;index"`
	ProcedureTemplate   ProcedureTemplate `json:"procedure_template" gorm:"foreignKey:ProcedureTemplateID"`

	// Procedure specific details for this treatment plan
	ToothNumber   string  `json:"tooth_number" gorm:"size:10"`
	Surface       string  `json:"surface" gorm:"size:50"`
	Notes         string  `json:"notes" gorm:"type:text"`
	EstimatedCost float64 `json:"estimated_cost" gorm:"type:decimal(10,2)"`
	Sequence      int     `json:"sequence" gorm:"default:1"` // order in treatment plan

	Status string `json:"status" gorm:"size:50;default:'planned'"` // planned, completed, skipped

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Helper methods for PatientDiagnosis
func (pd *PatientDiagnosis) IsActive() bool {
	return pd.Status == "active"
}

func (pd *PatientDiagnosis) IsResolved() bool {
	return pd.Status == "resolved"
}

func (pd *PatientDiagnosis) MarkAsResolved() {
	pd.Status = "resolved"
	now := time.Now()
	pd.ResolvedAt = &now
	pd.UpdatedAt = now
}

// Helper methods for PatientTreatmentPlan
func (ptp *PatientTreatmentPlan) IsActive() bool {
	return ptp.Status == "active"
}

func (ptp *PatientTreatmentPlan) IsCompleted() bool {
	return ptp.Status == "completed"
}

func (ptp *PatientTreatmentPlan) MarkAsCompleted() {
	ptp.Status = "completed"
	now := time.Now()
	ptp.CompletedAt = &now
	ptp.UpdatedAt = now
}

func (ptp *PatientTreatmentPlan) GetProgressPercentage() float64 {
	if ptp.EstimatedVisits == 0 {
		return 0
	}
	return (float64(ptp.CompletedVisits) / float64(ptp.EstimatedVisits)) * 100
}

func (ptp *PatientTreatmentPlan) GetRemainingVisits() int {
	return ptp.EstimatedVisits - ptp.CompletedVisits
}
