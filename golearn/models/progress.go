package models

import "gorm.io/gorm"

// Progress records when a student completes a specific lesson.
type Progress struct {
	gorm.Model
	UserID   uint `gorm:"not null;uniqueIndex:idx_user_lesson" json:"user_id"`
	LessonID uint `gorm:"not null;uniqueIndex:idx_user_lesson" json:"lesson_id"`
	CourseID uint `gorm:"not null" json:"course_id"`
}
