package handlers

import (
	"net/http"

	"golearn/database"
	"golearn/models"

	"github.com/gin-gonic/gin"
)

// CompleteLesson godoc
// @Summary      Mark a lesson as completed
// @Tags         progress
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /lessons/{id}/complete [post]
func CompleteLesson(c *gin.Context) {
	lessonID := c.Param("id")
	userID, _ := c.Get("user_id")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		return
	}

	// Prevent duplicate completion records
	var existing models.Progress
	if database.DB.Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(&existing).Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "lesson already completed"})
		return
	}

	progress := models.Progress{
		UserID:   userID.(uint),
		LessonID: lesson.ID,
		CourseID: lesson.CourseID,
	}
	database.DB.Create(&progress)
	c.JSON(http.StatusOK, gin.H{"message": "lesson marked as completed"})
}

// ProgressItem is the per-course summary returned by GetMyProgress.
type ProgressItem struct {
	CourseID         uint    `json:"course_id"`
	CourseTitle      string  `json:"course_title"`
	TotalLessons     int64   `json:"total_lessons"`
	CompletedLessons int64   `json:"completed_lessons"`
	Percent          float64 `json:"percent"`
}

// GetMyProgress godoc
// @Summary      Get current user's progress across all courses
// @Tags         progress
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   ProgressItem
// @Router       /my/progress [get]
func GetMyProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Find all courses in which the user has at least one completed lesson
	var progresses []models.Progress
	database.DB.Where("user_id = ?", userID).Find(&progresses)

	// Aggregate by course
	courseMap := map[uint]int64{}
	for _, p := range progresses {
		courseMap[p.CourseID]++
	}

	var result []ProgressItem
	for courseID, completed := range courseMap {
		var course models.Course
		if err := database.DB.First(&course, courseID).Error; err != nil {
			continue
		}

		var total int64
		database.DB.Model(&models.Lesson{}).Where("course_id = ?", courseID).Count(&total)

		pct := 0.0
		if total > 0 {
			pct = float64(completed) / float64(total) * 100
		}

		result = append(result, ProgressItem{
			CourseID:         courseID,
			CourseTitle:      course.Title,
			TotalLessons:     total,
			CompletedLessons: completed,
			Percent:          pct,
		})
	}

	c.JSON(http.StatusOK, result)
}
