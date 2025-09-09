package models

import (
	"time"

	"gorm.io/gorm"
)

type InquirySource string
type ConversionType string

const (
	SourceWebsite     InquirySource = "website"
	SourcePhone       InquirySource = "phone"
	SourceWalkIn      InquirySource = "walk_in"
	SourceReferral    InquirySource = "referral"
	SourceSocialMedia InquirySource = "social_media"
	SourceAdvertising InquirySource = "advertising"
)

const (
	ConversionInquiry     ConversionType = "inquiry"
	ConversionAppointment ConversionType = "appointment"
	ConversionTreatment   ConversionType = "treatment"
	ConversionPayment     ConversionType = "payment"
)

type Inquiry struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Source        InquirySource  `json:"source" gorm:"type:varchar(50);not null"`
	ContactName   string         `json:"contact_name" gorm:"size:200"`
	ContactPhone  string         `json:"contact_phone" gorm:"size:50"`
	ContactEmail  string         `json:"contact_email" gorm:"size:100"`
	Message       string         `json:"message" gorm:"type:text"`
	ServiceInterest string       `json:"service_interest" gorm:"size:200"`
	
	// Conversion tracking
	IsConverted   bool           `json:"is_converted" gorm:"default:false"`
	ConvertedAt   *time.Time     `json:"converted_at"`
	PatientID     *uint          `json:"patient_id" gorm:"index"`
	Patient       *Patient       `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	
	// Clinic association
	ClinicID      uint           `json:"clinic_id" gorm:"not null;index"`
	Clinic        Clinic         `json:"clinic" gorm:"foreignKey:ClinicID"`
	BranchID      *uint          `json:"branch_id" gorm:"index"`
	Branch        *Branch        `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	
	// Follow-up
	FollowUpDate  *time.Time     `json:"follow_up_date"`
	FollowUpNotes string         `json:"follow_up_notes" gorm:"type:text"`
	AssignedToID  *uint          `json:"assigned_to_id" gorm:"index"`
	AssignedTo    *User          `json:"assigned_to,omitempty" gorm:"foreignKey:AssignedToID"`
	
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type DailySales struct {
	ID                    uint           `json:"id" gorm:"primarykey"`
	Date                  time.Time      `json:"date" gorm:"type:date;not null;index"`
	
	// Appointment metrics
	TotalAppointments     int            `json:"total_appointments"`
	CompletedAppointments int            `json:"completed_appointments"`
	CancelledAppointments int            `json:"cancelled_appointments"`
	NoShowAppointments    int            `json:"no_show_appointments"`
	
	// Revenue metrics
	TotalRevenue          float64        `json:"total_revenue" gorm:"type:decimal(12,2)"`
	PaidRevenue          float64        `json:"paid_revenue" gorm:"type:decimal(12,2)"`
	PendingRevenue       float64        `json:"pending_revenue" gorm:"type:decimal(12,2)"`
	
	// Patient metrics
	NewPatients          int            `json:"new_patients"`
	ReturningPatients    int            `json:"returning_patients"`
	
	// Inquiry and conversion metrics
	TotalInquiries       int            `json:"total_inquiries"`
	ConvertedInquiries   int            `json:"converted_inquiries"`
	ConversionRate       float64        `json:"conversion_rate" gorm:"type:decimal(5,2)"`
	
	// Clinic association
	ClinicID             uint           `json:"clinic_id" gorm:"not null;index"`
	Clinic               Clinic         `json:"clinic" gorm:"foreignKey:ClinicID"`
	
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

type PatientAnalytics struct {
	ID                    uint           `json:"id" gorm:"primarykey"`
	PatientID             uint           `json:"patient_id" gorm:"not null;index"`
	Patient               Patient        `json:"patient" gorm:"foreignKey:PatientID"`
	
	// Visit patterns
	TotalVisits           int            `json:"total_visits"`
	LastVisitDate         *time.Time     `json:"last_visit_date"`
	AverageTimeBetweenVisits int         `json:"avg_time_between_visits"` // in days
	
	// Financial metrics
	TotalSpent            float64        `json:"total_spent" gorm:"type:decimal(12,2)"`
	AverageSpentPerVisit  float64        `json:"avg_spent_per_visit" gorm:"type:decimal(10,2)"`
	OutstandingBalance    float64        `json:"outstanding_balance" gorm:"type:decimal(10,2)"`
	
	// Attendance patterns
	NoShowCount           int            `json:"no_show_count"`
	LateArrivalCount      int            `json:"late_arrival_count"`
	CancellationCount     int            `json:"cancellation_count"`
	ReliabilityScore      float64        `json:"reliability_score" gorm:"type:decimal(3,2)"` // 0-1 score
	
	// Treatment patterns
	MostCommonTreatments  string         `json:"most_common_treatments" gorm:"type:text"`
	TreatmentCompliance   float64        `json:"treatment_compliance" gorm:"type:decimal(3,2)"`
	
	// Communication preferences
	PreferredContactMethod string        `json:"preferred_contact_method" gorm:"size:50"`
	ResponsivenessScore   float64        `json:"responsiveness_score" gorm:"type:decimal(3,2)"`
	
	LastCalculated        time.Time      `json:"last_calculated"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// Methods for Inquiry
func (i *Inquiry) MarkAsConverted(patientID uint) {
	i.IsConverted = true
	now := time.Now()
	i.ConvertedAt = &now
	i.PatientID = &patientID
}

// Methods for DailySales
func (ds *DailySales) CalculateConversionRate() {
	if ds.TotalInquiries > 0 {
		ds.ConversionRate = float64(ds.ConvertedInquiries) / float64(ds.TotalInquiries) * 100
	} else {
		ds.ConversionRate = 0
	}
}

// Methods for PatientAnalytics
func (pa *PatientAnalytics) CalculateReliabilityScore() {
	if pa.TotalVisits == 0 {
		pa.ReliabilityScore = 1.0
		return
	}
	
	// Calculate reliability based on no-shows, late arrivals, and cancellations
	totalProblems := pa.NoShowCount + pa.LateArrivalCount + pa.CancellationCount
	problemRate := float64(totalProblems) / float64(pa.TotalVisits)
	pa.ReliabilityScore = 1.0 - problemRate
	
	// Ensure score stays between 0 and 1
	if pa.ReliabilityScore < 0 {
		pa.ReliabilityScore = 0
	}
}