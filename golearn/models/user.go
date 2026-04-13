package models

import "gorm.io/gorm"

// Role constants used for RBAC.
const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

// User represents a platform member (student or teacher).
type User struct {
	gorm.Model
	Name     string `gorm:"not null"          json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null"          json:"-"`
	Role     string `gorm:"default:student"   json:"role"`

	// Associations
	Courses []Course `gorm:"foreignKey:TeacherID" json:"courses,omitempty"`
}
