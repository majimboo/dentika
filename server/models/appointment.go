package models

import (
	"time"

	"gorm.io/gorm"
)

type AppointmentStatus string

const (
	StatusScheduled   AppointmentStatus = "scheduled"
	StatusConfirmed   AppointmentStatus = "confirmed"
	StatusInProgress  AppointmentStatus = "in_progress"
	StatusCompleted   AppointmentStatus = "completed"
	StatusCancelled   AppointmentStatus = "cancelled"
	StatusNoShow      AppointmentStatus = "no_show"
	StatusRescheduled AppointmentStatus = "rescheduled"
)

type Appointment struct {
	ID          uint              `json:"id" gorm:"primarykey"`
	Title       string            `json:"title" gorm:"size:200"`
	Description string            `json:"description" gorm:"type:text"`
	StartTime   time.Time         `json:"start_time" gorm:"not null;index"`
	EndTime     time.Time         `json:"end_time" gorm:"not null"`
	Duration    int               `json:"duration"` // in minutes
	Status      AppointmentStatus `json:"status" gorm:"type:varchar(20);default:'scheduled'"`

	// Patient arrival tracking
	PatientArrived bool       `json:"patient_arrived" gorm:"default:false"`
	ArrivalTime    *time.Time `json:"arrival_time"`
	IsLate         bool       `json:"is_late" gorm:"default:false"`

	// Notes and follow-up
	PreAppointmentNotes  string     `json:"pre_appointment_notes" gorm:"type:text"`
	PostAppointmentNotes string     `json:"post_appointment_notes" gorm:"type:text"`
	NextAppointmentDate  *time.Time `json:"next_appointment_date"`

	// Cost and payment
	EstimatedCost float64 `json:"estimated_cost" gorm:"type:decimal(10,2)"`
	ActualCost    float64 `json:"actual_cost" gorm:"type:decimal(10,2)"`
	IsPaid        bool    `json:"is_paid" gorm:"default:false"`

	// Foreign Keys
	PatientID uint    `json:"patient_id" gorm:"not null;index"`
	Patient   Patient `json:"patient" gorm:"foreignKey:PatientID"`

	DoctorID uint `json:"doctor_id" gorm:"not null;index"`
	Doctor   User `json:"doctor" gorm:"foreignKey:DoctorID"`

	BranchID uint   `json:"branch_id" gorm:"not null;index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Relationships
	Procedures []AppointmentProcedure `json:"procedures,omitempty" gorm:"foreignKey:AppointmentID"`
	Diagnoses  []AppointmentDiagnosis `json:"diagnoses,omitempty" gorm:"foreignKey:AppointmentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AppointmentReminder struct {
	ID            uint        `json:"id" gorm:"primarykey"`
	AppointmentID uint        `json:"appointment_id" gorm:"not null;index"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`
	ReminderType  string      `json:"reminder_type" gorm:"size:50"` // email, sms, call
	ReminderTime  time.Time   `json:"reminder_time" gorm:"not null"`
	IsSent        bool        `json:"is_sent" gorm:"default:false"`
	SentAt        *time.Time  `json:"sent_at"`
	Message       string      `json:"message" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Appointment) IsUpcoming() bool {
	return a.StartTime.After(time.Now()) && a.Status == StatusScheduled
}

func (a *Appointment) IsToday() bool {
	now := time.Now()
	appointmentDate := a.StartTime
	return appointmentDate.Year() == now.Year() &&
		appointmentDate.YearDay() == now.YearDay()
}

func (a *Appointment) GetTimeUntilAppointment() time.Duration {
	return time.Until(a.StartTime)
}

func (a *Appointment) ShouldShowCountdown() bool {
	timeUntil := a.GetTimeUntilAppointment()
	return timeUntil > 0 && timeUntil <= 30*time.Minute && a.Status == StatusScheduled
}

func (a *Appointment) MarkPatientArrived() {
	a.PatientArrived = true
	now := time.Now()
	a.ArrivalTime = &now

	// Check if patient is late (more than 10 minutes after scheduled time)
	if now.After(a.StartTime.Add(10 * time.Minute)) {
		a.IsLate = true
	}
}
