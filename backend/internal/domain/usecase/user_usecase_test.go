package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"


	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/w0ikid/study-platform/internal/app/config"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
	"gitlab.com/w0ikid/study-platform/internal/domain/services"
	"gitlab.com/w0ikid/study-platform/internal/dto"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
)

// Mock для UserService
type MockUserService struct {
	mock.Mock
	services.UserServiceInterface
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) SearchUsers(ctx context.Context, name string) ([]*models.User, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		input := &dto.CreateUserInput{
			Username: "testuser",
			Name:     "Test",
			Surname:  "User",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "student",
		}

		expectedUser := &models.User{
			ID:        1,
			Username:  input.Username,
			Name:      input.Name,
			Surname:   input.Surname,
			Email:     input.Email,
			Password:  input.Password, // В реальности здесь был бы хеш
			Role:      input.Role,
			CreatedAt: time.Now(),
		}

		mockService.On("CreateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == input.Username &&
				u.Name == input.Name &&
				u.Email == input.Email
		})).Return(expectedUser, nil).Once()

		user, err := useCase.CreateUser(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, input.Username, user.Username)
		assert.Equal(t, input.Email, user.Email)
		mockService.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Отсутствует обязательное поле Email
		input := &dto.CreateUserInput{
			Username: "testuser",
			Name:     "Test",
			Surname:  "User",
			Password: "password123",
			Role:     "student",
		}

		user, err := useCase.CreateUser(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "validation failed")
		mockService.AssertNotCalled(t, "CreateUser")
	})

	t.Run("Service Error", func(t *testing.T) {
		input := &dto.CreateUserInput{
			Username: "testuser",
			Name:     "Test",
			Surname:  "User",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "student",
		}

		serviceError := errors.New("database error")
		mockService.On("CreateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == input.Username &&
				u.Name == input.Name &&
				u.Email == input.Email
		})).Return(nil, serviceError).Once()

		user, err := useCase.CreateUser(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, serviceError, err)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := 1
		expectedUser := &models.User{
			ID:       userID,
			Username: "testuser",
			Name:     "Test",
			Surname:  "User",
			Email:    "test@example.com",
			Role:     "student",
		}

		mockService.On("GetUser", ctx, userID).Return(expectedUser, nil).Once()

		user, err := useCase.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := 999
		mockService.On("GetUser", ctx, userID).Return(nil, errors.New("user not found")).Once()

		user, err := useCase.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockService.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &models.User{
			ID:       1,
			Username: "testuser",
			Name:     "Test",
			Surname:  "User",
			Email:    email,
			Role:     "student",
		}

		mockService.On("GetUserByEmail", ctx, email).Return(expectedUser, nil).Once()

		user, err := useCase.GetUserByEmail(ctx, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		email := "nonexistent@example.com"
		mockService.On("GetUserByEmail", ctx, email).Return(nil, errors.New("user not found")).Once()

		user, err := useCase.GetUserByEmail(ctx, email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockService.AssertExpectations(t)
	})
}

func TestGetUserByUsername(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		username := "testuser"
		expectedUser := &models.User{
			ID:       1,
			Username: username,
			Name:     "Test",
			Surname:  "User",
			Email:    "test@example.com",
			Role:     "student",
		}

		mockService.On("GetUserByUsername", ctx, username).Return(expectedUser, nil).Once()

		user, err := useCase.GetUserByUsername(ctx, username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		username := "nonexistent"
		mockService.On("GetUserByUsername", ctx, username).Return(nil, errors.New("user not found")).Once()

		user, err := useCase.GetUserByUsername(ctx, username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		expectedUser := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    email,
			Role:     "student",
		}

		mockService.On("Login", ctx, email, password).Return(expectedUser, nil).Once()

		user, token, err := useCase.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"

		mockService.On("Login", ctx, email, password).Return(nil, errors.New("invalid credentials")).Once()

		user, token, err := useCase.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockService.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := 1
		mockService.On("DeleteUser", ctx, userID).Return(nil).Once()

		err := useCase.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := 999
		mockService.On("DeleteUser", ctx, userID).Return(errors.New("user not found")).Once()

		err := useCase.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockService.AssertExpectations(t)
	})
}

func TestSearchUsers(t *testing.T) {
	mockService := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:       "test-secret",
			ExpiredHours: 24,
		},
	}
	useCase := usecase.NewUserUseCase(mockService, cfg)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		searchName := "Test"
		expectedUsers := []*models.User{
			{
				ID:       1,
				Username: "testuser1",
				Name:     "Test",
				Surname:  "User1",
				Email:    "test1@example.com",
				Role:     "student",
			},
			{
				ID:       2,
				Username: "testuser2",
				Name:     "Test",
				Surname:  "User2",
				Email:    "test2@example.com",
				Role:     "student",
			},
		}

		mockService.On("SearchUsers", ctx, searchName).Return(expectedUsers, nil).Once()

		users, err := useCase.SearchUsers(ctx, searchName)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		assert.Equal(t, expectedUsers[0].ID, users[0].ID)
		assert.Equal(t, expectedUsers[1].ID, users[1].ID)
		mockService.AssertExpectations(t)
	})

	t.Run("No Users Found", func(t *testing.T) {
		searchName := "Nonexistent"
		mockService.On("SearchUsers", ctx, searchName).Return([]*models.User{}, nil).Once()

		users, err := useCase.SearchUsers(ctx, searchName)

		assert.NoError(t, err)
		assert.Empty(t, users)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		searchName := "Error"
		mockService.On("SearchUsers", ctx, searchName).Return(nil, errors.New("database error")).Once()

		users, err := useCase.SearchUsers(ctx, searchName)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "database error")
		mockService.AssertExpectations(t)
	})
}

func TestToUserResponses(t *testing.T) {
	t.Run("Convert Users to Responses", func(t *testing.T) {
		users := []*models.User{
			{
				ID:        1,
				Username:  "user1",
				Email:     "user1@example.com",
				Role:      "student",
				Level:     2,
				Xp:        150,
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				Username:  "user2",
				Email:     "user2@example.com",
				Role:      "teacher",
				Level:     5,
				Xp:        750,
				CreatedAt: time.Now(),
			},
		}

		responses := usecase.ToUserResponses(users)

		assert.Len(t, responses, 2)
		assert.Equal(t, users[0].ID, responses[0].ID)
		assert.Equal(t, users[0].Username, responses[0].Username)
		assert.Equal(t, users[0].Email, responses[0].Email)
		assert.Equal(t, users[0].Role, responses[0].Role)
		assert.Equal(t, users[0].Level, responses[0].Level)
		assert.Equal(t, users[0].Xp, responses[0].Xp)
		assert.Equal(t, users[0].CreatedAt, responses[0].CreatedAt)

		assert.Equal(t, users[1].ID, responses[1].ID)
		assert.Equal(t, users[1].Username, responses[1].Username)
		assert.Equal(t, users[1].Email, responses[1].Email)
		assert.Equal(t, users[1].Role, responses[1].Role)
		assert.Equal(t, users[1].Level, responses[1].Level)
		assert.Equal(t, users[1].Xp, responses[1].Xp)
		assert.Equal(t, users[1].CreatedAt, responses[1].CreatedAt)
	})

	t.Run("Empty User List", func(t *testing.T) {
		users := []*models.User{}
		responses := usecase.ToUserResponses(users)
		assert.Empty(t, responses)
	})

	t.Run("Nil User List", func(t *testing.T) {
		var users []*models.User = nil
		responses := usecase.ToUserResponses(users)
		assert.Empty(t, responses)
	})
}