package models

import (
	"time"

	"gorm.io/gorm"
)

type Clinic struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name" gorm:"size:200;not null"`
	Code     string `json:"code" gorm:"size:50;uniqueIndex"` // Unique clinic code for URLs
	Address  string `json:"address" gorm:"size:500"`
	Phone    string `json:"phone" gorm:"size:50"`
	Email    string `json:"email" gorm:"size:100"`
	Website  string `json:"website" gorm:"size:200"`
	Logo     string `json:"logo" gorm:"size:500"`
	IsActive bool   `json:"is_active" gorm:"default:true"`

	// Relationships
	Branches []Branch  `json:"branches,omitempty" gorm:"foreignKey:ClinicID"`
	Staff    []User    `json:"staff,omitempty" gorm:"foreignKey:ClinicID"`
	Patients []Patient `json:"patients,omitempty" gorm:"foreignKey:ClinicID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Branch struct {
	ID           uint   `json:"id" gorm:"primarykey"`
	Name         string `json:"name" gorm:"size:200;not null"`
	Address      string `json:"address" gorm:"size:500"`
	Phone        string `json:"phone" gorm:"size:50"`
	IsMainBranch bool   `json:"is_main_branch" gorm:"default:false"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`

	// Foreign Keys
	ClinicID uint   `json:"clinic_id" gorm:"not null;index"`
	Clinic   Clinic `json:"clinic" gorm:"foreignKey:ClinicID"`

	// Relationships
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:BranchID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Clinic) GetMainBranch() *Branch {
	for _, branch := range c.Branches {
		if branch.IsMainBranch {
			return &branch
		}
	}
	return nil
}

func (c *Clinic) GetActiveBranches() []Branch {
	var activeBranches []Branch
	for _, branch := range c.Branches {
		if branch.IsActive {
			activeBranches = append(activeBranches, branch)
		}
	}
	return activeBranches
}
