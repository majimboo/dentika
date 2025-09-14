package models

import (
	"time"

	"gorm.io/gorm"
)

type DailySales struct {
	ID   uint      `json:"id" gorm:"primarykey"`
	Date time.Time `json:"date" gorm:"type:date;not null;index"`

	// Appointment metrics
	TotalAppointments     int `json:"total_appointments"`
	CompletedAppointments int `json:"completed_appointments"`
	CancelledAppointments int `json:"cancelled_appointments"`
	NoShowAppointments    int `json:"no_show_appointments"`

	// Revenue metrics
	TotalRevenue   float64 `json:"total_revenue" gorm:"type:decimal(12,2)"`
	PaidRevenue    float64 `json:"paid_revenue" gorm:"type:decimal(12,2)"`
	PendingRevenue float64 `json:"pending_revenue" gorm:"type:decimal(12,2)"`

	// Patient metrics
	NewPatients       int `json:"new_patients"`
	ReturningPatients int `json:"returning_patients"`

	// Conversion metrics
	TotalInquiries     int     `json:"total_inquiries"`
	ConvertedInquiries int     `json:"converted_inquiries"`
	ConversionRate     float64 `json:"conversion_rate" gorm:"type:decimal(5,2)"`

	// Clinic association
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Methods for DailySales
func (ds *DailySales) CalculateConversionRate() {
	if ds.TotalInquiries > 0 {
		ds.ConversionRate = float64(ds.ConvertedInquiries) / float64(ds.TotalInquiries) * 100
	} else {
		ds.ConversionRate = 0
	}
}
