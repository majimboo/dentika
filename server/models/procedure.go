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
