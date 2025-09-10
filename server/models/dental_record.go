package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ToothType string
type ToothCondition string

const (
	ToothTypePrimary   ToothType = "primary"   // baby teeth
	ToothTypePermanent ToothType = "permanent" // adult teeth
)

const (
	ConditionHealthy   ToothCondition = "healthy"
	ConditionDecay     ToothCondition = "decay"
	ConditionFilled    ToothCondition = "filled"
	ConditionCrowned   ToothCondition = "crowned"
	ConditionRootCanal ToothCondition = "root_canal"
	ConditionExtracted ToothCondition = "extracted"
	ConditionImplant   ToothCondition = "implant"
	ConditionBridge    ToothCondition = "bridge"
	ConditionMissing   ToothCondition = "missing"
)

type DentalRecord struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	PatientID uint    `json:"patient_id" gorm:"not null;index"`
	Patient   Patient `json:"patient" gorm:"foreignKey:PatientID"`

	// Clinic scoping for multi-tenancy
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	RecordType ToothType `json:"record_type" gorm:"type:varchar(20);default:'permanent'"` // primary or permanent teeth
	IsActive   bool      `json:"is_active" gorm:"default:true"`

	// JSON field to store all tooth conditions
	TeethData string `json:"teeth_data" gorm:"type:text"` // JSON array of tooth objects

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SurfaceCondition struct {
	Surface   string         `json:"surface"`   // mesial, distal, occlusal, buccal, lingual, incisal
	Condition ToothCondition `json:"condition"` // specific condition for this surface
}

type ToothData struct {
	ToothNumber       string             `json:"tooth_number"`       // e.g., "1", "2", ..., "32" for permanent, "A", "B", ... for primary
	Position          string             `json:"position"`           // quadrant and position info
	Condition         ToothCondition     `json:"condition"`          // overall tooth condition
	Surfaces          []string           `json:"surfaces"`           // affected surfaces (for backward compatibility)
	SurfaceConditions []SurfaceCondition `json:"surface_conditions"` // detailed conditions per surface
	Notes             string             `json:"notes"`
	LastUpdated       time.Time          `json:"last_updated"`
	UpdatedBy         uint               `json:"updated_by"` // user ID who made the update
}

type DentalRecordHistory struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	DentalRecordID uint         `json:"dental_record_id" gorm:"not null;index"`
	DentalRecord   DentalRecord `json:"dental_record" gorm:"foreignKey:DentalRecordID"`

	ToothNumber       string         `json:"tooth_number" gorm:"size:10;not null"`
	PreviousCondition ToothCondition `json:"previous_condition" gorm:"type:varchar(50)"`
	NewCondition      ToothCondition `json:"new_condition" gorm:"type:varchar(50)"`
	ChangeReason      string         `json:"change_reason" gorm:"type:text"`
	AppointmentID     *uint          `json:"appointment_id" gorm:"index"`
	Appointment       *Appointment   `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`

	ChangedByID uint `json:"changed_by_id" gorm:"not null;index"`
	ChangedBy   User `json:"changed_by" gorm:"foreignKey:ChangedByID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Methods for DentalRecord
func (dr *DentalRecord) GetTeethData() ([]ToothData, error) {
	var teethData []ToothData
	if dr.TeethData == "" {
		return dr.initializeTeethData(), nil
	}

	err := json.Unmarshal([]byte(dr.TeethData), &teethData)
	if err != nil {
		return nil, err
	}
	return teethData, nil
}

func (dr *DentalRecord) SetTeethData(teethData []ToothData) error {
	data, err := json.Marshal(teethData)
	if err != nil {
		return err
	}
	dr.TeethData = string(data)
	return nil
}

func (dr *DentalRecord) initializeTeethData() []ToothData {
	var teethData []ToothData

	if dr.RecordType == ToothTypePermanent {
		// Initialize 32 permanent teeth using FDI notation
		quadrants := []struct {
			quadrantNum int
			startPos    int
			endPos      int
		}{
			{1, 1, 8}, // Upper right (11-18)
			{2, 1, 8}, // Upper left (21-28)
			{3, 1, 8}, // Lower left (31-38)
			{4, 1, 8}, // Lower right (41-48)
		}

		for _, quad := range quadrants {
			for pos := quad.startPos; pos <= quad.endPos; pos++ {
				toothNum := quad.quadrantNum*10 + pos
				tooth := ToothData{
					ToothNumber:       fmt.Sprintf("%d", toothNum),
					Position:          dr.getToothPosition(toothNum),
					Condition:         ConditionHealthy,
					Surfaces:          []string{},
					SurfaceConditions: []SurfaceCondition{},
					Notes:             "",
					LastUpdated:       time.Now(),
				}
				teethData = append(teethData, tooth)
			}
		}
	} else {
		// Initialize 20 primary teeth (A-T)
		letters := "ABCDEFGHIJKLMNOPQRST"
		for i, letter := range letters {
			tooth := ToothData{
				ToothNumber:       string(letter),
				Position:          dr.getPrimaryToothPosition(i + 1),
				Condition:         ConditionHealthy,
				Surfaces:          []string{},
				SurfaceConditions: []SurfaceCondition{},
				Notes:             "",
				LastUpdated:       time.Now(),
			}
			teethData = append(teethData, tooth)
		}
	}

	return teethData
}

func (dr *DentalRecord) getToothPosition(toothNumber int) string {
	quadrants := map[int]string{
		1: "Upper Right", 2: "Upper Left",
		3: "Lower Left", 4: "Lower Right",
	}

	// For FDI notation, extract quadrant from first digit
	quadrant := toothNumber / 10
	position := toothNumber % 10

	return fmt.Sprintf("%s - %d", quadrants[quadrant], position)
}

func (dr *DentalRecord) getPrimaryToothPosition(toothNumber int) string {
	quadrants := map[int]string{
		1: "Upper Right", 2: "Upper Left",
		3: "Lower Left", 4: "Lower Right",
	}

	quadrant := ((toothNumber - 1) / 5) + 1
	position := ((toothNumber - 1) % 5) + 1

	return fmt.Sprintf("%s - %d", quadrants[quadrant], position)
}

func (dr *DentalRecord) UpdateToothCondition(toothNumber string, newCondition ToothCondition, surfaces []string, surfaceConditions []SurfaceCondition, notes string, updatedBy uint) error {
	teethData, err := dr.GetTeethData()
	if err != nil {
		return err
	}

	for i := range teethData {
		if teethData[i].ToothNumber == toothNumber {
			teethData[i].Condition = newCondition
			teethData[i].Surfaces = surfaces
			teethData[i].SurfaceConditions = surfaceConditions
			teethData[i].Notes = notes
			teethData[i].LastUpdated = time.Now()
			teethData[i].UpdatedBy = updatedBy
			break
		}
	}

	return dr.SetTeethData(teethData)
}

// Helper method to get surface condition
func (td *ToothData) GetSurfaceCondition(surface string) ToothCondition {
	for _, sc := range td.SurfaceConditions {
		if sc.Surface == surface {
			return sc.Condition
		}
	}
	return ConditionHealthy
}

// Helper method to set surface condition
func (td *ToothData) SetSurfaceCondition(surface string, condition ToothCondition) {
	// Remove existing condition for this surface
	for i, sc := range td.SurfaceConditions {
		if sc.Surface == surface {
			td.SurfaceConditions = append(td.SurfaceConditions[:i], td.SurfaceConditions[i+1:]...)
			break
		}
	}

	// Add new condition if not healthy
	if condition != ConditionHealthy {
		td.SurfaceConditions = append(td.SurfaceConditions, SurfaceCondition{
			Surface:   surface,
			Condition: condition,
		})
	}
}
