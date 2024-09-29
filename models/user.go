// models/user.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email          string   `gorm:"uniqueIndex;not null" json:"email"`
    Username       string   `gorm:"uniqueIndex;not null" json:"username"`
    Password       string   `gorm:"not null" json:"password,omitempty"`
    PhoneNumber    string   `json:"phone_number"`
    ProfilePicture string   `json:"profile_picture"` // Field baru untuk URL profile picture
    PackageID      *uint    `json:"package_id,omitempty"`
    Package        Package  `json:"package,omitempty"`
}
