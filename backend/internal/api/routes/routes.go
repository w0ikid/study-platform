package routes

import (
	"github.com/gin-gonic/gin"

	"gitlab.com/w0ikid/study-platform/internal/api/handlers"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
	"gitlab.com/w0ikid/study-platform/internal/api/middlewares"
	"gitlab.com/w0ikid/study-platform/internal/app/config"
)

func SetupRoutes(r *gin.Engine, userUseCase *usecase.UserUseCase, courseUseCase *usecase.CourseUseCase, lessonUseCase *usecase.LessonUseCase ,enrollment *usecase.EnrollmentUseCase , lessonProgressUseCase *usecase.LessonProgressUseCase , certificateUseCase *usecase.CertificateUseCase, cfg *config.Config) {
	userHandler := handlers.NewUserHandler(userUseCase)
	courseHandler := handlers.NewCourseHandler(courseUseCase)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollment)
	lessonHandler := handlers.NewLessonHandler(lessonUseCase)
	lessonProgressHandler := handlers.NewLessonProgressHandler(lessonProgressUseCase)
	certificateHandler := handlers.NewCertificateHandler(certificateUseCase)
	// Middlewares
	authMiddleware := middlewares.AuthMiddleware(cfg.JWT)
	enrollmentMiddleware := middlewares.EnrollmentMiddleware(enrollment)
	// enrollmentByLesson := middlewares.EnrollmentByLessonMiddleware(lessonUseCase, enrollment)
	api := r.Group("/api")
	{	
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.CreateUser)
			// auth.GET("/me", userHandler.GetMe)
		}
		// Users
		users := api.Group("/users")
		{
			users.GET("/:username", authMiddleware, userHandler.GetUserByUsername)
			users.DELETE("/:id", authMiddleware, middlewares.RoleMiddleware("admin"), userHandler.DeleteUser)
			users.GET("/", authMiddleware, userHandler.SearchUsers)	
		}
		// Courses
		courses := api.Group("/courses")
		{
			courses.POST("/", authMiddleware, middlewares.RoleMiddleware("admin", "teacher"), courseHandler.CreateCourse)
			courses.GET("/:id", authMiddleware, courseHandler.GetCourse)
			courses.GET("/", authMiddleware, courseHandler.GetAllCourses)
			

			// enrollments
			courses.POST("/:id/enroll", authMiddleware, enrollmentHandler.CreateEnrollment)
			// courses.GET("/:course_id/enrollments", authMiddleware, enrollmentHandler.GetEnrollmentsByCourse)
			
			// courses.PUT("/:id", courseHandler.UpdateCourse)
			courses.DELETE("/:id", authMiddleware, middlewares.RoleMiddleware("admin", "teacher"), courseHandler.Delete)

			// lessons
			courses.POST("/:id/lessons", authMiddleware, middlewares.RoleMiddleware("admin", "teacher"), lessonHandler.CreateLesson)
			courses.GET("/:id/lessons", authMiddleware, enrollmentMiddleware, lessonHandler.GetLessonsByCourse)
			courses.POST("/:id/lessons/:lesson_id/complete", authMiddleware, enrollmentMiddleware, lessonProgressHandler.CompleteLesson)
			// lesson progress
			courses.GET("/:id/progress", authMiddleware, enrollmentMiddleware, lessonProgressHandler.GetCourseProgress)

			// certificates
			// courses.GET("/:id/certificate", authMiddleware, certificateHandler.GenerateCertificate)


		}
		enrollment := api.Group("/enrollment")
		{
			enrollment.GET("/:id", authMiddleware, enrollmentHandler.GetEnrollment)
			enrollment.GET("/", authMiddleware, enrollmentHandler.GetAllEnrollment)
			enrollment.DELETE("/:id", authMiddleware, middlewares.RoleMiddleware("admin"), enrollmentHandler.Delete)
		}
		certificates := api.Group("/certificates")
		{
			certificates.GET("/course/:course_id", authMiddleware, certificateHandler.GenerateCertificate)
		}
		// lessons := api.Group("lessons")
		// {
		// 	lessons.POST("/:id/complete", authMiddleware, enrollmentByLesson, lessonProgressHandler.CompleteLesson)
		// }
		
	}
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

