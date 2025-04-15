package controllers

import (
	"net/http"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/review_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	rs *review_service.ReviewService
}

func NewReviewController(rs *review_service.ReviewService) *ReviewController {
	return &ReviewController{rs: rs}
}

// @Summary Создание нового отзыва (от пользователя к пользователю)
// @Description Создает отзыв по email пользователя
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body dto.ReviewToCreateDto true "Review to create"
// @Success 201 {object} exceptions.Error_
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /reviews [post]
func (rc *ReviewController) CreateReview(ctx *gin.Context) {
	var review dto.ReviewToCreateDto
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exc := review.Validate(); exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

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

	reviewerEmail := payload["email"].(string)

	if exc := rc.rs.CreateReview(ctx.Request.Context(), reviewerEmail, &review); exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(201, gin.H{"success": true})
}

// @Summary Получение отзывов для пользователя
// @Description Возвращает список отзывов для указанного пользователя (по username)
// @Tags reviews
// @Produce json
// @Param username path string true "Username"
// @Success 200 {array} dto.ReviewToGetDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 404 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /users/{username}/reviews [get]
func (rc *ReviewController) GetReviewsForUser(ctx *gin.Context) {
	targetUsername := ctx.Param("username")

	reviews, exc := rc.rs.GetReviewsForUser(ctx.Request.Context(), targetUsername, 0, 0)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*reviews) == 0 {
		data := make([]string, 0)
		ctx.JSON(404, data)
		return
	}

	ctx.JSON(200, reviews)
}
