package models

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole defines the roles a user can have in the system.
// Possible values are: super_admin, clinic_owner, doctor, secretary, assistant.
type UserRole string

const (
	SuperAdmin UserRole = "super_admin"
	ClinicOwner UserRole = "clinic_owner"
	Doctor UserRole = "doctor"
	Secretary UserRole = "secretary"
	Assistant UserRole = "assistant"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
	FirstName string         `json:"first_name" gorm:"size:100"`
	LastName  string         `json:"last_name" gorm:"size:100"`
	Gender    string         `json:"gender" gorm:"size:20"`
	Avatar    string         `json:"avatar" gorm:"size:500"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"type:varchar(20);default:'secretary'"`
	ClinicID  *uint          `json:"clinic_id" gorm:"index"`
	Clinic    *Clinic        `json:"clinic,omitempty" gorm:"foreignKey:ClinicID"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type AuthToken struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Token       string         `json:"token" gorm:"uniqueIndex;not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	SessionData string         `json:"session_data" gorm:"type:text"`
	ExpiresAt   time.Time      `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (t *AuthToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

func (u *User) GetDisplayName() string {
	fullName := u.GetFullName()
	if fullName == u.Username {
		return u.Username
	}
	return fullName
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == SuperAdmin
}

func (u *User) CanAccessClinic(clinicID uint) bool {
	if u.IsSuperAdmin() {
		return true
	}
	return u.ClinicID != nil && *u.ClinicID == clinicID
}

func (u *User) HasRole(roles ...UserRole) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}
