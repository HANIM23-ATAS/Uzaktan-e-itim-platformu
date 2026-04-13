package models

import "gorm.io/gorm"

// Quiz is a set of questions attached to a lesson.
type Quiz struct {
	gorm.Model
	Title    string     `gorm:"not null" json:"title"`
	LessonID uint       `gorm:"uniqueIndex;not null" json:"lesson_id"`
	Questions []Question `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
}

// Question is a multiple-choice item within a quiz.
type Question struct {
	gorm.Model
	Text    string `gorm:"not null" json:"text"`
	OptionA string `json:"option_a"`
	OptionB string `json:"option_b"`
	OptionC string `json:"option_c"`
	OptionD string `json:"option_d"`
	Correct string `gorm:"not null" json:"correct"` // "a", "b", "c", or "d"
	QuizID  uint   `gorm:"not null" json:"quiz_id"`
}

// QuizResult stores the outcome of a student's quiz attempt.
type QuizResult struct {
	gorm.Model
	UserID  uint    `gorm:"not null" json:"user_id"`
	QuizID  uint    `gorm:"not null" json:"quiz_id"`
	Score   int     `json:"score"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}
