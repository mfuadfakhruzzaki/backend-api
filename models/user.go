package models

import (
	"time"
)

type User struct {
    ID              uint        `gorm:"primarykey" json:"id"`
    CreatedAt       time.Time   `json:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at"`
    DeletedAt       *time.Time  `json:"deleted_at,omitempty"`

    Email           string      `gorm:"uniqueIndex;not null" json:"email"`
    Username        string      `gorm:"uniqueIndex;not null" json:"username"`
    Password        string      `gorm:"not null" json:"password,omitempty"`
    PhoneNumber     string      `json:"phone_number"`
    ProfilePicture  string      `json:"profile_picture"`
    PackageID       *uint       `json:"package_id,omitempty"`
    Package         Package     `json:"package,omitempty"`
    EmailVerified   bool        `gorm:"default:false" json:"email_verified"`
    VerificationCode string     `gorm:"size:6" json:"-"`
}
