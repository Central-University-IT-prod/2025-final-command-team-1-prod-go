package controllers

import (
	"strconv"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/place_service"
	"github.com/gin-gonic/gin"
)

type PlaceController struct {
	ps place_service.PlaceService
}

func NewPlaceController(s *place_service.PlaceService) *PlaceController {
	return &PlaceController{
		ps: *s,
	}
}

// CreatePlace создает новое место
// @Summary Создать место
// @Description Создает новое место (доступно только администраторам)
// @Tags places
// @Accept json
// @Produce json
// @Param place body dto.CreatePlaceDto true "Данные нового места"
// @Success 201 {object} dto.PlaceDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /places [post]
func (c *PlaceController) CreatePlace(ctx *gin.Context) {
	var place dto.CreatePlaceDto
	if err := ctx.ShouldBindBodyWithJSON(&place); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdPlace, err := c.ps.CreatePlace(ctx, &place)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	ctx.JSON(201, &createdPlace)
}

// GetPlaces возвращает все места
// @Summary Получить все места
// @Description Возвращает все места
// @Tags places
// @Accept json
// @Produce json
// @Success 200 {array} dto.PlaceDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /places [get]
func (c *PlaceController) GetPlaces(ctx *gin.Context) {
	places, err := c.ps.GetPlaces(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	if len(*places) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}

	ctx.JSON(200, places)
}

// DeletePlace удаляет место
// @Summary Удалить место
// @Description Удаляет место по ID
// @Tags places
// @Accept json
// @Produce json
// @Param id path int true "ID места"
// @Success 204 {object} nil
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 404 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /places/{id} [delete]
func (c *PlaceController) DeletePlace(ctx *gin.Context) {
	id := ctx.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	exc := c.ps.DeletePlace(ctx, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}
	ctx.JSON(204, nil)
}
