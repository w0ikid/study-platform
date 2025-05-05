package middlewares

import (
	_"log"
	"net/http"
	"strings"

	"gitlab.com/w0ikid/study-platform/internal/app/config"
	"gitlab.com/w0ikid/study-platform/pkg/auth"

	"strconv"

	"github.com/gin-gonic/gin"
	
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

func AuthMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение токена из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Проверка формата токена
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Проверка токена
		claims, err := auth.ValidateJWT(parts[1], jwtConfig.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Установка данных пользователя в контекст
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("userRole")
		
		// Проверка наличия роли в списке разрешенных
		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 
func EnrollmentMiddleware(enrollmentUseCase *usecase.EnrollmentUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID := c.GetInt("userID")
		courseID, err := strconv.Atoi(c.Param("id"))
		
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			c.Abort()
			return
		}
		
		userRole := c.GetString("userRole")
		if userRole == "admin" || userRole == "teacher" {
			c.Next()
			return
		}
		
		// Проверяем, зачислен ли пользователь на курс
		enrolled, err := enrollmentUseCase.IsUserEnrolled(ctx, userID, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		
		if !enrolled {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: you are not enrolled in this course"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// func EnrollmentByLessonMiddleware(
//     lessonUseCase	*usecase.LessonUseCase,
//     enrollmentUseCase *usecase.EnrollmentUseCase,
// ) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         ctx := c.Request.Context()
//         userID := c.GetInt("userID")
//         lessonID, err := strconv.Atoi(c.Param("id"))
//         if err != nil {
//             c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
//             c.Abort()
//             return
//         }

//         userRole := c.GetString("userRole")
//         if userRole == "admin" || userRole == "teacher" {
//             c.Next()
//             return
//         }

//         lesson, err := lessonUseCase.GetLessonByID(ctx, lessonID)
//         if err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Lesson not found"})
//             c.Abort()
//             return
//         }

//         enrolled, err := enrollmentUseCase.IsUserEnrolled(ctx, userID, lesson.CourseID)
//         if err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//             c.Abort()
//             return
//         }

//         if !enrolled {
//             c.JSON(http.StatusForbidden, gin.H{"error": "You are not enrolled in this course"})
//             c.Abort()
//             return
//         }

//         c.Next()
//     }
// }
