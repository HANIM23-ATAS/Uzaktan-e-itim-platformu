package handlers

import (
	"net/http"
	"strconv"

	"golearn/database"
	"golearn/models"

	"github.com/gin-gonic/gin"
)

// CourseInput defines the request body for creating/updating a course.
type CourseInput struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ListCourses godoc
// @Summary      List courses
// @Description  Returns paginated list of courses with optional category filter and sort
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Param        page     query int    false "Page number"
// @Param        limit    query int    false "Items per page"
// @Param        category query string false "Filter by category"
// @Param        sort     query string false "Sort field (created_at|title)"
// @Success      200  {object}  map[string]interface{}
// @Router       /courses [get]
func ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	query := database.DB.Model(&models.Course{}).Preload("Teacher")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	var courses []models.Course
	query.Order(sort).Offset(offset).Limit(limit).Find(&courses)

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetCourse godoc
// @Summary      Get course by ID
// @Description  Returns a single course with lessons and teacher info
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200  {object}  models.Course
// @Failure      404  {object}  map[string]string
// @Router       /courses/{id} [get]
func GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := database.DB.Preload("Teacher").Preload("Lessons").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// CreateCourse godoc
// @Summary      Create a course (teacher only)
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CourseInput true "Course data"
// @Success      201  {object}  models.Course
// @Failure      400  {object}  map[string]string
// @Router       /courses [post]
func CreateCourse(c *gin.Context) {
	var input CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	course := models.Course{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		TeacherID:   userID.(uint),
	}

	database.DB.Create(&course)
	database.DB.Preload("Teacher").First(&course, course.ID)
	c.JSON(http.StatusCreated, course)
}

// UpdateCourse godoc
// @Summary      Update a course (owner teacher only)
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int         true "Course ID"
// @Param        body body CourseInput true "Updated course data"
// @Success      200  {object}  models.Course
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /courses/{id} [put]
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this course"})
		return
	}

	var input CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&course).Updates(models.Course{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
	})
	c.JSON(http.StatusOK, course)
}

// DeleteCourse godoc
// @Summary      Delete a course (owner teacher only)
// @Tags         courses
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /courses/{id} [delete]
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this course"})
		return
	}

	database.DB.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}
