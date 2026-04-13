package models

import "gorm.io/gorm"

// Lesson is a single unit of content inside a course.
type Lesson struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title"`
	Content  string `json:"content"`
	VideoURL string `json:"video_url"`
	Order    int    `gorm:"default:0" json:"order"`
	CourseID uint   `gorm:"not null" json:"course_id"`

	// Association
	Quiz *Quiz `gorm:"foreignKey:LessonID" json:"quiz,omitempty"`
}
