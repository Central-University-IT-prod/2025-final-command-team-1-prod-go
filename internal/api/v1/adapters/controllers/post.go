package controllers

import (
	"strconv"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/core/application/services/post_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	ps post_service.PostService
}

func NewPostController(s *post_service.PostService) *PostController {
	return &PostController{
		ps: *s,
	}
}

// @Summary Создать объявление
// @Description Создать новое объявление.
// @Tags posts
// @Accept json
// @Produce json
// @Param post body dto.CreatePostDto true "Post details"
// @Success 201 {object} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts [post]
func (c *PostController) CreatePost(ctx *gin.Context) {
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

	var post dto.CreatePostDto
	if err := ctx.ShouldBindBodyWithJSON(&post); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdPost, err := c.ps.CreatePost(ctx, email, &post)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(201, &createdPost)
	go c.ps.SetPostSummary(ctx, createdPost.ID, createdPost.Author, createdPost.Title)
	go c.ps.SetPostQuote(ctx, createdPost.ID, createdPost.Author, createdPost.Title)
}

// @Summary Получить объявление по ID.
// @Description Возвращает объявление по его айди. Если такого не существует, то возвращает 404.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 404 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id} [get]
func (c *PostController) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")

	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

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

	post, exc := c.ps.GetPost(ctx, idParsed, email)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(200, &post)
}

// @Summary Получить объявления пользователя
// @Description Получить все объявления пользователя.
// @Tags posts
// @Accept json
// @Produce json
// @Param status query string false "Post status" Enums(available, booked, taken, all) default(all)
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.PostToGetDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/my [get]
func (c *PostController) GetMyPosts(ctx *gin.Context) {
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

	status := ctx.Query("status")
	if status == "" {
		status = "all"
	}

	if status != "available" && status != "booked" && status != "taken" && status != "all" {
		ctx.JSON(400, gin.H{"error": "Invalid status"})
		return
	}

	limit, _ := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseUint(ctx.Query("offset"), 10, 64)

	posts, exc := c.ps.GetMyPosts(ctx, email, status, uint(limit), uint(offset))
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*posts) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}
	ctx.JSON(200, &posts)
}

// @Summary Добавить пост в избранное
// @Description Добавляет пост в избранное. Если пост уже в избранном, то ничего не происходит.
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} exceptions.Error_
// @Failure 404 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id}/favorites [put]
func (c *PostController) AddFavorite(ctx *gin.Context) {
	id := ctx.Param("id")

	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

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
	exc = c.ps.AddFavorite(ctx, email, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.JSON(200, &dto.OKStatus)
}

// @Summary Удалить пост из избранного
// @Description Удаляет пост из избранного. Если пост не добавлен в избранное, то ничего не происходит.
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 204
// @Failure 400 {object} exceptions.Error_
// @Failure 404 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/{id}/favorites [delete]
func (c *PostController) DeleteFavorite(ctx *gin.Context) {
	id := ctx.Param("id")

	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

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
	exc = c.ps.DeleteFavorite(ctx, email, idParsed)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	ctx.Status(204)
}

// @Summary Получить все доступные посты
// @Description Получить все доступные посты. Есть возможность отфильтровать.
// @Tags posts
// @Accept json
// @Produce json
// @Param genre query string false "Post genre"
// @Param condition query string false "Post condition"
// @Param publicationYear query string false "Post publication year"
// @Param placeId query string false "Post place id"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/available [get]
func (c *PostController) GetAllAvailablePosts(ctx *gin.Context) {
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

	limit, _ := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseUint(ctx.Query("offset"), 10, 64)

	placeIDStr := ctx.Query("placeId")
	placeID, err := strconv.ParseInt(placeIDStr, 10, 64)
	if err != nil {
		placeID = 0 // or handle the error as needed
	}
	options := repositories.PostFilterOptions{
		Genre:           ctx.Query("genre"),
		Condition:       ctx.Query("condition"),
		PublicationYear: ctx.Query("publicationYear"),
		PlaceID:         placeID,
	}

	posts, exc := c.ps.GetAllAvailablePosts(ctx, email, options, uint(limit), uint(offset))
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*posts) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}

	ctx.JSON(200, &posts)
}

// @Summary Получить все избранные посты
// @Description Получить все избранные посты пользователя.
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/favorites [get]
func (c *PostController) GetAllFavourites(ctx *gin.Context) {
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

	limit, _ := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseUint(ctx.Query("offset"), 10, 64)

	posts, exc := c.ps.GetAllFavourites(ctx, email, uint(limit), uint(offset))
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*posts) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}

	ctx.JSON(200, &posts)
}

// @Summary Поиск постов по названию или автору
// @Description Поиск постов по названию или автору.
// @Tags posts
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/search [get]
func (r *PostController) SearchByTitleOrAuthorOrGenre(ctx *gin.Context) {
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

	query := ctx.Query("query")
	limit, _ := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseUint(ctx.Query("offset"), 10, 64)

	posts, exc := r.ps.SearchByTitleOrAuthorOrGenre(ctx, query, uint(limit), uint(offset), email)
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*posts) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}

	ctx.JSON(200, &posts)
}

// @Summary Получить все забронированные посты
// @Description Получить все посты, забронированные пользователем.
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.PostDto
// @Failure 400 {object} exceptions.Error_
// @Failure 401 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Failure 503 {object} exceptions.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Router /posts/booked [get]
func (c *PostController) GetAllMyBooked(ctx *gin.Context) {
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

	limit, _ := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseUint(ctx.Query("offset"), 10, 64)

	posts, exc := c.ps.GetAllMyBooked(ctx, email, uint(limit), uint(offset))
	if exc != nil {
		ctx.JSON(int(exc.StatusCode), exc)
		return
	}

	if len(*posts) == 0 {
		data := make([]string, 0)
		ctx.JSON(200, data)
		return
	}

	ctx.JSON(200, &posts)
}

// @Summary Добавление изображения к объявлению
// @Schemes
// @Description Загружает изображение для рекламной кампании
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID поста"
// @Param image formData file true "Файл изображения"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.Response
// @Failure 400 {object} exceptions.Error_
// @Failure 415 {object} exceptions.Error_
// @Failure 500 {object} exceptions.Error_
// @Security BearerAuth
// @Router /posts/{postID}/image [post]
func (c *PostController) AddImage(ctx *gin.Context) {
	id := ctx.Param("id")

	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

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

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(int(exceptions.ErrNotAllFields.StatusCode), exceptions.ErrNotAllFields)
		return
	}
	defer file.Close()
	go c.ps.AddImage(ctx, idParsed, header, file, email)
	ctx.JSON(200, gin.H{"message": "ok"})
}
