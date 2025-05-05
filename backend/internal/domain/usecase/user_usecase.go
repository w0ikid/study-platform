package usecase

import (
	"context"
	"gitlab.com/w0ikid/study-platform/internal/app/config"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
	"gitlab.com/w0ikid/study-platform/internal/dto"
	"gitlab.com/w0ikid/study-platform/pkg/auth"
	"github.com/go-playground/validator/v10"
	"fmt"
)

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	SearchUsers(ctx context.Context, name string) ([]*models.User, error)
}

type UserUseCase struct {
	userService services.UserServiceInterface
	jwtConfig 	config.JWTConfig
}

func NewUserUseCase(userService services.UserServiceInterface, cfg *config.Config) *UserUseCase {
	return &UserUseCase{userService: userService, jwtConfig: cfg.JWT }
}

func (u *UserUseCase) CreateUser(ctx context.Context, input *dto.CreateUserInput) (*models.User, error) {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	user := &models.User{
		Username: input.Username,
		Name: 	  input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Password: input.Password,
		Role: 	  input.Role,
	}
	return u.userService.CreateUser(ctx, user)
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	// Дополнительная логика перед получением пользователя (например, проверка прав доступа)
	return u.userService.GetUser(ctx, id)
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Дополнительная логика перед получением пользователя (например, проверка прав доступа)
	return u.userService.GetUserByEmail(ctx, email)
}


func (u *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	// Дополнительная логика перед получением пользователя (например, проверка прав доступа)
	return u.userService.GetUserByUsername(ctx, username)
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (*models.User, string, error) {
    user, err := u.userService.Login(ctx, email, password)
    if err != nil {
        return nil, "", err
    }

    // Генерация JWT
    token, err := auth.GenerateJWT(user.ID, user.Role, u.jwtConfig.Secret, u.jwtConfig.ExpiredHours)
    if err != nil {
        return nil, "", fmt.Errorf("failed to generate token: %w", err)
    }

    return user, token, nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id int) error {
	// Дополнительная логика перед удалением пользователя (например, проверка прав доступа)
	// Вызов сервиса для удаления пользователя
	return u.userService.DeleteUser(ctx, id)
}

func (u *UserUseCase) SearchUsers(ctx context.Context, name string) ([]*models.User, error) {
	// Дополнительная логика перед поиском пользователей (например, фильтрация)
	return u.userService.SearchUsers(ctx, name)
}

func ToUserResponses(users []*models.User) []models.UserResponse {
	var responses []models.UserResponse
	for _, u := range users {
		responses = append(responses, models.UserResponse{
			ID:       	u.ID,
			Username: 	u.Username,
			Email:    	u.Email,
			Role:     	u.Role,
			Level:  	u.Level,
			Xp:        	u.Xp,
			CreatedAt: 	u.CreatedAt,

		})
	}
	return responses
}