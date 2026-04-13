package middleware

import (
	"net/http"

	"golearn/models"

	"github.com/gin-gonic/gin"
)

// TeacherOnly blocks any request whose role claim is not "teacher".
func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != models.RoleTeacher {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "teacher access required"})
			return
		}
		c.Next()
	}
}
