package handlers

import (
	_ "log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"strings"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
	"gitlab.com/w0ikid/study-platform/internal/dto"
	
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase

}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser godoc
// @Summary      Register new user
// @Description  Register a new user account with username, email, password and role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserInput  true  "User Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string  "Validation error"
// @Failure      409   {object}  map[string]string  "Email already taken"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /auth/register [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request dto.CreateUserInput

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	user, err := h.userUseCase.CreateUser(ctx, &request)
	
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "email is already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email is already taken"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"level":    user.Level,
		"xp":       user.Xp,
	})
}

// @Summary      Get user by ID
// @Description  Retrieve a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := c.Request.Context()

	user, err := h.userUseCase.GetUserByID(ctx, int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name": 	user.Name,
		"surname": 	user.Surname,
		"email":    user.Email,
		"role":     user.Role,
		"level":    user.Level,
		"xp":       user.Xp,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	ctx := c.Request.Context()

	user, err := h.userUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"level":    user.Level,
		"xp":       user.Xp,
		"created_at": user.CreatedAt,
	})
}
// @Summary      Get user by username
// @Description  Returns user info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username path string true "Username"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/{username} [get]
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	ctx := c.Request.Context()

	user, err := h.userUseCase.GetUserByUsername(ctx, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"level":    user.Level,
		"xp":       user.Xp,
		"created_at": user.CreatedAt,
	})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginUserInput  true  "Login credentials"
// @Success      200          {object}  map[string]string   "JWT token"
// @Failure      400          {object}  map[string]string   "Invalid input"
// @Failure      401          {object}  map[string]string   "Invalid credentials"
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var input dto.LoginUserInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	_, token, err := h.userUseCase.Login(ctx, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// // jwtToken
	// token, err := auth.GenerateJWT(user.ID, user.Role, h.jwtConfig.Secret, h.jwtConfig.ExpiredHours)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	// 	return
	// }

	// refreshToken, err := auth.GenerateJWT(user.ID, user.Role, h.jwtConfig.RefreshSecret, h.jwtConfig.RefreshExpiresHours)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
	// 	return
	// }

	// Возвращаем ответ
	// c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}


// DeleteUser обрабатывает удаление пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем контекст из запроса
	ctx := c.Request.Context()

	// Используем UseCase для удаления пользователя
	err = h.userUseCase.DeleteUser(ctx, int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	name := c.Query("name")

	// Получаем контекст из запроса
	ctx := c.Request.Context()

	// Используем UseCase для поиска пользователей
	users, err := h.userUseCase.SearchUsers(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// Преобразуем пользователей в нужный формат
	safeUsers := usecase.ToUserResponses(users)

	// Возвращаем ответ
	c.JSON(http.StatusOK, gin.H{
		"users": safeUsers,
	})
}

// refresh
// func (h *UserHandler) Refresh(c *gin.Context) {
// 	refreshToken, err := c.Cookie("refresh_token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
// 		return
// 	}

// 	claims, err := auth.ValidateJWT(refreshToken, h.jwtConfig.RefreshSecret)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 		return
// 	}

// 	// Выдаем новый access token
// 	newToken, err := auth.GenerateJWT(claims.UserID, claims.Role, h.jwtConfig.Secret, h.jwtConfig.ExpiredHours)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": newToken})
// }