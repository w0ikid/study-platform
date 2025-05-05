package services

import (
	"context" // добавляем импорт context
	"gitlab.com/w0ikid/study-platform/internal/domain/repositories"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	SearchUsers(ctx context.Context, name string) ([]*models.User, error)
	UpdateXpAndLevel(ctx context.Context, user *models.User) error
}

type UserService struct {
	repo repositories.UserRepositoryInterface
}

func NewUserService(repo repositories.UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: repo}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)


	// Прокидываем ctx в репозиторий
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser получает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

// GetUserByEmail получает пользователя по email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.FindByUsername(ctx, username)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// compare the password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) SearchUsers(ctx context.Context, name string) ([]*models.User, error) {
	return s.repo.FindByNameLike(ctx, name)
}

func (s *UserService) UpdateXpAndLevel(ctx context.Context, user *models.User) error {
	return s.repo.UpdateXpAndLevel(ctx, user)
}