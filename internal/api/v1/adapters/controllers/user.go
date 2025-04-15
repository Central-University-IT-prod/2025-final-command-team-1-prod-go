package controllers

import (
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/user_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	us user_service.UserService
}

func NewUserController(s *user_service.UserService) *UserController {
	return &UserController{
		us: *s,
	}
}

// @BasePath /api/v1

// @Summary Бинд токена для уведомлений
// @Schemes
// @Description Биндит к юзеру токен мобильного устройства
// @Tags user
// @Accept json
// @Produce json
// @Param token body dto.BindTokenDto true "Token data"
// @Success 201 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Router /users/bind_token [post]
func (c *UserController) BindPushToken(ctx *gin.Context) {
	var token dto.BindTokenDto
	if err := ctx.ShouldBindBodyWithJSON(&token); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authToken, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*authToken)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	email := payload["email"].(string)
	if err := c.us.BindPushToken(ctx, email, token.Token); err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(201, gin.H{"success": true})
}

// @Summary Создание пользователя
// @Schemes
// @Description Creates new user and returns it
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.CreateUserDto true "User data"
// @Success 201 {object} dto.GetUserDto
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Failure 400 {object} exceptions.Error_
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user dto.CreateUserDto
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := c.us.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	userToReturn := utils.ExcludeUserCredentials(createdUser)
	ctx.JSON(201, &userToReturn)
}

// @Summary Получить профиль пользователя (по нику)
// @Schemes
// @Description Returns user profile by username (requires JWT in "Bearer" header)
// @Tags user
// @Produce json
// @Param username path string true "Username"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.UserWithRatingDto
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Security BearerAuth
// @Router /users/{username} [get]
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.us.GetUserByUsernameWithRating(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	// userToReturn := utils.ExcludeUserCredentials(user)

	ctx.JSON(200, &user)
}

// @Summary Получить профиль пользователя (свой)
// @Schemes
// @Description Returns user profile (requires JWT in "Bearer" header)
// @Tags user
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.UserWithRatingDto
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Security BearerAuth
// @Router /users/me [get]
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	user, err := c.us.GetUserByUsernameWithRating(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	// userToReturn := utils.ExcludeUserCredentials(user)

	ctx.JSON(200, &user)
}

// @Summary Обновить профиль.
// @Schemes
// @Description Updates user profile and returns it (requires JWT in "Bearer" header)
// @Tags user
// @Produce json
// @Param user body dto.UpdateUserDto true "User data"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.GetUserDto
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Security BearerAuth
// @Router /users/me [patch]
func (c *UserController) UpdateUserProfile(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	email := payload["email"].(string)

	var updateData dto.UpdateUserDto
	if err := ctx.ShouldBindBodyWithJSON(&updateData); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := c.us.UpdateUserByEmail(ctx, email, updateData)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	userToReturn := utils.ExcludeUserCredentials(updatedUser)

	ctx.JSON(200, &userToReturn)
}
