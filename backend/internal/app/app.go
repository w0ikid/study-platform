package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/app/config"
	"gitlab.com/w0ikid/study-platform/internal/app/connections"
	"gitlab.com/w0ikid/study-platform/internal/app/start"
	_ "gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

func Run(configFile string) error {

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return err
	}

	// Инициализация соединений
	conn, err := connections.NewConnections(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	// Автомиграция моделей
	if err := autoMigrate(conn.DB); err != nil {
		return err
	}

	fmt.Printf("%s - DBname\n", cfg.DB.DBName)

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(conn.DB)
	courseRepo := repositories.NewCourseRepository(conn.DB)
	certificateRepo := repositories.NewCertificateRepository(conn.DB)
	enrollmentRepo := repositories.NewEnrollmentRepository(conn.DB)
	lessonRepo := repositories.NewLessonRepository(conn.DB)
	lessonProgressRepo := repositories.NewLessonProgressRepository(conn.DB)
	// Инициализация сервисов
	userService := services.NewUserService(userRepo)
	courseService := services.NewCourseService(courseRepo)
	certificateService := services.NewCertificateService(certificateRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo)
	lessonService := services.NewLessonService(lessonRepo)
	lessonProgressService := services.NewLessonProgressService(lessonProgressRepo)
	// Инициализация usecase
	userUseCase := usecase.NewUserUseCase(userService, cfg)
	courseUseCase := usecase.NewCourseUseCase(courseService)
	enrollmentUseCase := usecase.NewEnrollmentUseCase(enrollmentService, courseService)
	lessonUseCase := usecase.NewLessonUseCase(lessonService, enrollmentService, courseService)
	lessonProgressUseCase := usecase.NewLessonProgressUseCase(lessonProgressService, lessonService, enrollmentService, courseService, userService)
	certificateUseCase := usecase.NewCertificateUseCase(certificateService, enrollmentService, userService, courseService)
	// Запуск HTTP сервера
	start.HTTP(cfg, userUseCase, courseUseCase, lessonUseCase, enrollmentUseCase, lessonProgressUseCase, certificateUseCase)

	return nil
}

// autoMigrate выполняет автоматическую миграцию для всех моделей
func autoMigrate(conn *pgx.Conn) error {
	ctx := context.Background()
	fmt.Println("Starting auto migration...")
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			name VARCHAR(52),
			surname VARCHAR(52),
			email TEXT NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role TEXT NOT NULL,
			xp INT DEFAULT 0,
			level INT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS courses (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			image_url TEXT,
			teacher_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'active',  -- Может быть 'active', 'inactive',
			CONSTRAINT fk_teacher FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS lessons (
			id SERIAL PRIMARY KEY,
			course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			video_url TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS lesson_progress (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			lesson_id INT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
			course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE, -- Для удобства связки с курсом
			is_completed BOOLEAN DEFAULT FALSE, -- Завершен ли урок
			completed_at TIMESTAMP, -- Время завершения урока (если завершен)
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uq_lesson_progress_user_lesson UNIQUE(user_id, lesson_id), -- Уникальность комбинации user_id и lesson_id
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
			CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS enrollments (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			status TEXT DEFAULT 'active', -- Персональный статус студента, completed - выполнено
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
			CONSTRAINT uq_enrollment_user_course UNIQUE(user_id, course_id)  -- уникальность комбинации user_id и course_id
		);`,
		`CREATE TABLE IF NOT EXISTS certificates (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uq_certificate_user_course UNIQUE(user_id, course_id)  -- Гарантия одного сертификата на курс
		);`,
	}

	for i, query := range queries {
		fmt.Printf("Executing query %d...\n", i+1)
		_, err := conn.Exec(ctx, query)
		if err != nil {
			fmt.Printf("Error executing query %d: %v\n", i+1, err)
			return fmt.Errorf("error executing query: %v", err)
		}
		fmt.Printf("Query %d executed successfully\n", i+1)
	}

	fmt.Println("Auto migration completed successfully")

	return nil
}