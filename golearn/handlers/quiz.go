package handlers

import (
	"net/http"

	"golearn/database"
	"golearn/models"

	"github.com/gin-gonic/gin"
)

// QuestionInput defines a single question within a quiz creation request.
type QuestionInput struct {
	Text    string `json:"text"    binding:"required"`
	OptionA string `json:"option_a"`
	OptionB string `json:"option_b"`
	OptionC string `json:"option_c"`
	OptionD string `json:"option_d"`
	Correct string `json:"correct" binding:"required"`
}

// QuizInput is the request body for creating a quiz.
type QuizInput struct {
	Title     string          `json:"title"     binding:"required"`
	Questions []QuestionInput `json:"questions" binding:"required,min=1"`
}

// SubmitInput is the request body for submitting quiz answers.
type SubmitInput struct {
	// Answers maps question_id (as string) to answered letter ("a"-"d").
	Answers map[uint]string `json:"answers" binding:"required"`
}

// GetQuiz godoc
// @Summary      Get quiz for a lesson
// @Tags         quiz
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Success      200  {object}  models.Quiz
// @Failure      404  {object}  map[string]string
// @Router       /lessons/{id}/quiz [get]
func GetQuiz(c *gin.Context) {
	lessonID := c.Param("id")
	var quiz models.Quiz
	if err := database.DB.Preload("Questions").Where("lesson_id = ?", lessonID).First(&quiz).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found for this lesson"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// CreateQuiz godoc
// @Summary      Create quiz for a lesson (teacher/owner only)
// @Tags         quiz
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int       true "Lesson ID"
// @Param        body body QuizInput true "Quiz data"
// @Success      201  {object}  models.Quiz
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /lessons/{id}/quiz [post]
func CreateQuiz(c *gin.Context) {
	lessonID := c.Param("id")

	// Verify lesson exists and fetch its course for ownership check
	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		return
	}

	var course models.Course
	database.DB.First(&course, lesson.CourseID)

	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this course"})
		return
	}

	// Only one quiz per lesson
	var existing models.Quiz
	if database.DB.Where("lesson_id = ?", lessonID).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "quiz already exists for this lesson"})
		return
	}

	var input QuizInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz := models.Quiz{
		Title:    input.Title,
		LessonID: lesson.ID,
	}
	database.DB.Create(&quiz)

	// Bulk-insert questions
	for _, q := range input.Questions {
		question := models.Question{
			Text:    q.Text,
			OptionA: q.OptionA,
			OptionB: q.OptionB,
			OptionC: q.OptionC,
			OptionD: q.OptionD,
			Correct: q.Correct,
			QuizID:  quiz.ID,
		}
		database.DB.Create(&question)
	}

	database.DB.Preload("Questions").First(&quiz, quiz.ID)
	c.JSON(http.StatusCreated, quiz)
}

// SubmitQuiz godoc
// @Summary      Submit quiz answers
// @Tags         quiz
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int         true "Quiz ID"
// @Param        body body SubmitInput true "Answers map"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /quiz/{id}/submit [post]
func SubmitQuiz(c *gin.Context) {
	quizID := c.Param("id")

	var quiz models.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	var input SubmitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score := 0
	total := len(quiz.Questions)
	for _, q := range quiz.Questions {
		if ans, ok := input.Answers[q.ID]; ok && ans == q.Correct {
			score++
		}
	}

	percent := 0.0
	if total > 0 {
		percent = float64(score) / float64(total) * 100
	}

	userID, _ := c.Get("user_id")
	result := models.QuizResult{
		UserID:  userID.(uint),
		QuizID:  quiz.ID,
		Score:   score,
		Total:   total,
		Percent: percent,
	}
	database.DB.Create(&result)

	message := "Keep practicing!"
	if percent >= 80 {
		message = "Great job!"
	} else if percent >= 50 {
		message = "Good effort!"
	}

	c.JSON(http.StatusOK, gin.H{
		"score":   score,
		"total":   total,
		"percent": percent,
		"message": message,
	})
}
