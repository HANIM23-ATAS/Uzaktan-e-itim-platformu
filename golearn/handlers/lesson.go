package handlers

import (
	"net/http"

	"golearn/database"
	"golearn/models"

	"github.com/gin-gonic/gin"
)

// LessonInput defines the request body for creating a lesson.
type LessonInput struct {
	Title    string `json:"title"    binding:"required"`
	Content  string `json:"content"`
	VideoURL string `json:"video_url"`
	Order    int    `json:"order"`
}

// ListLessons godoc
// @Summary      List lessons for a course
// @Tags         lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200  {array}   models.Lesson
// @Failure      404  {object}  map[string]string
// @Router       /courses/{id}/lessons [get]
func ListLessons(c *gin.Context) {
	courseID := c.Param("id")
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	var lessons []models.Lesson
	database.DB.Where("course_id = ?", courseID).Order("`order` asc").Find(&lessons)
	c.JSON(http.StatusOK, lessons)
}

// CreateLesson godoc
// @Summary      Add a lesson to a course (owner teacher only)
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int         true "Course ID"
// @Param        body body LessonInput true "Lesson data"
// @Success      201  {object}  models.Lesson
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /courses/{id}/lessons [post]
func CreateLesson(c *gin.Context) {
	courseID := c.Param("id")
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	// Only the course owner can add lessons
	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this course"})
		return
	}

	var input LessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson := models.Lesson{
		Title:    input.Title,
		Content:  input.Content,
		VideoURL: input.VideoURL,
		Order:    input.Order,
		CourseID: course.ID,
	}
	database.DB.Create(&lesson)
	c.JSON(http.StatusCreated, lesson)
}
