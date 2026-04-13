package models

import "gorm.io/gorm"

// Course represents a learning course created by a teacher.
type Course struct {
	gorm.Model
	Title       string `gorm:"not null"  json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	TeacherID   uint   `gorm:"not null"  json:"teacher_id"`

	// Associations
	Teacher User     `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:CourseID"  json:"lessons,omitempty"`
}
