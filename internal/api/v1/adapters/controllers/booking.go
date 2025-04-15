package controllers

import (
	"strconv"

	"example.com/m/internal/api/v1/core/application/services/booking_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bs booking_service.BookingService
}

func NewBookingController(bs *booking_service.BookingService) *BookingController {
	return &BookingController{bs: *bs}
}

// @Summary Забронировать книгу
// @Description Забронировать книгу по ID
// @Tags бронирования
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Success 201 {object} map[string]interface{} "success"
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id}/booking [post]
func (c *BookingController) BookBook(ctx *gin.Context) {
	token, exc := utils.ExtractTokenFromHeaders(ctx)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	payload, exc := utils.ExtractPayloadFromJWT(*token)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	email := payload["email"].(string)

	id := ctx.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	exc = c.bs.BookBook(ctx, email, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(201, gin.H{"success": true})
}

// @Summary Delete a booking
// @Description Delete a booking by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id}/booking [delete]
func (c *BookingController) DeleteBooking(ctx *gin.Context) {

	token, exc := utils.ExtractTokenFromHeaders(ctx)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	payload, exc := utils.ExtractPayloadFromJWT(*token)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	email := payload["email"].(string)

	id := ctx.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	exc = c.bs.DeleteBooking(ctx, email, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(200, gin.H{"success": true})
}

// @Summary Mark a book as taken
// @Description Mark a booking as taken by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id}/mark-taken [put]
func (c *BookingController) MarkAsTaken(ctx *gin.Context) {
	token, exc := utils.ExtractTokenFromHeaders(ctx)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	payload, exc := utils.ExtractPayloadFromJWT(*token)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	email := payload["email"].(string)

	id := ctx.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	exc = c.bs.MarkAsTaken(ctx, email, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(200, gin.H{"success": true})
}
