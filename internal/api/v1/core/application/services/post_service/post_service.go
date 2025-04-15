package post_service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/core/application/services/gpt_service"
	"example.com/m/internal/api/v1/core/application/services/place_service"
	"example.com/m/internal/api/v1/core/application/services/user_service"
)

type PostService struct {
	pr repositories.PostRepository
	ps place_service.PlaceService
	us user_service.UserService
	gs gpt_service.GPTService
}

func NewPostService(
	pr *repositories.PostRepository, ps *place_service.PlaceService,
	us *user_service.UserService, gs *gpt_service.GPTService,
) *PostService {
	return &PostService{pr: *pr, ps: *ps, us: *us, gs: *gs}
}

func (ps *PostService) CreatePost(ctx context.Context, userEmail string, p *dto.CreatePostDto) (*dto.PostDto, *exceptions.Error_) {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	if !*userExists {
		return nil, &exceptions.ErrUserNotFound
	}

	placeExists, exc := ps.ps.PlaceIsExists(ctx, p.PlaceID)
	if exc != nil {
		return nil, exc
	}
	if !*placeExists {
		return nil, &exceptions.ErrPlaceNotFound
	}

	postToCreate := dto.PostDtoWithoutId{
		Images:          p.Images,
		UserEmail:       userEmail,
		PlaceID:         p.PlaceID,
		Title:           p.Title,
		Description:     p.Description,
		Genre:           p.Genre,
		Author:          p.Author,
		PublicationYear: p.PublicationYear,
		Publisher:       p.Publisher,
		Condition:       p.Condition,
		Status:          "available",
		CreatedAt:       time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Cover:           p.Cover,
		PagesCount:      p.PagesCount,
	}

	id, err := ps.pr.Create(ctx, &postToCreate)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	post, err := ps.pr.Get(ctx, *id, userEmail)
	if err != nil {
		fmt.Println(err)
		return nil, &exceptions.ErrDatabaseError
	}

	return post, nil
}

func (ps *PostService) SetPostSummary(ctx context.Context, postID int64, author string, bookName string) *exceptions.Error_ {
	summary, err := ps.gs.GenerateBriefContent(ctx, bookName, author)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	er := ps.pr.SetSummary(ctx, postID, summary)
	if er != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}

func (ps *PostService) SetPostQuote(ctx context.Context, postID int64, author string, bookName string) *exceptions.Error_ {
	summary, err := ps.gs.GenerateQuote(ctx, bookName, author)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	er := ps.pr.SetQuote(ctx, postID, summary)
	if er != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}

func (ps *PostService) GetMyPosts(ctx context.Context, userEmail string, status string, limit, offset uint) (*[]dto.PostDto, *exceptions.Error_) {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return nil, exc
	}
	if !*userExists {
		return nil, &exceptions.ErrUserNotFound
	}

	posts, err := ps.pr.GetAllMyPosted(ctx, userEmail, status, limit, offset)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return posts, nil
}

func (ps *PostService) GetPost(ctx context.Context, id int64, userEmail string) (*dto.PostDto, *exceptions.Error_) {
	post, err := ps.pr.Get(ctx, id, userEmail)
	if err != nil {
		fmt.Println(err)
		return nil, &exceptions.ErrDatabaseError
	}

	if post == nil {
		return nil, &exceptions.PostNotFoudErr
	}

	return post, nil
}

func (ps *PostService) UserIsPostOwner(ctx context.Context, userEmail string, id int64) (*bool, *exceptions.Error_) {
	post, err := ps.pr.Get(ctx, id, userEmail)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	state := post == nil || post.UserEmail == userEmail
	return &state, nil
}

func (ps *PostService) DeletePost(ctx context.Context, userEmail string, id int64) *exceptions.Error_ {
	if userIsOwner, err := ps.UserIsPostOwner(ctx, userEmail, id); err != nil {
		return &exceptions.ErrDatabaseError
	} else if !*userIsOwner {
		return &exceptions.ErrUserIsNotOwner
	}

	err := ps.pr.Delete(ctx, id)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}

func (ps *PostService) GetAllAvailablePosts(ctx context.Context, userEmail string, options repositories.PostFilterOptions, limit, offset uint) (*[]dto.PostDto, *exceptions.Error_) {
	posts, err := ps.pr.GetAllAvailable(ctx, userEmail, options, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, &exceptions.ErrDatabaseError
	}

	return posts, nil
}

func (ps *PostService) UpdatePost(ctx context.Context, userEmail string, id int64, p *dto.UpdatePostDto) (*dto.PostDto, *exceptions.Error_) {
	if userIsOwner, err := ps.UserIsPostOwner(ctx, userEmail, id); err != nil {
		return nil, &exceptions.ErrDatabaseError
	} else if !*userIsOwner {
		return nil, &exceptions.ErrUserIsNotOwner
	}

	err := ps.pr.Update(ctx, id, p)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	post, err := ps.pr.Get(ctx, id, userEmail)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	return post, nil
}

func (ps *PostService) GetAllMyPosts(ctx context.Context, userEmail string, status string, limit, offset uint) (*[]dto.PostDto, *exceptions.Error_) {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return nil, exc
	}
	if !*userExists {
		return nil, &exceptions.ErrUserNotFound
	}

	posts, err := ps.pr.GetAllMyPosted(ctx, userEmail, status, limit, offset)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	return posts, nil
}

func (ps *PostService) DeleteFavorite(ctx context.Context, userEmail string, postID int64) *exceptions.Error_ {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return exc
	}
	if !*userExists {
		return &exceptions.ErrUserNotFound
	}
	post, err := ps.pr.Get(ctx, postID, userEmail)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	if post == nil {
		return &exceptions.PostNotFoudErr
	}
	err = ps.pr.DeleteFavorite(ctx, postID, userEmail)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}

func (ps *PostService) AddFavorite(ctx context.Context, userEmail string, postID int64) *exceptions.Error_ {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return exc
	}
	if !*userExists {
		return &exceptions.ErrUserNotFound
	}
	post, err := ps.pr.Get(ctx, postID, userEmail)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	if post == nil {
		return &exceptions.PostNotFoudErr
	}
	err = ps.pr.AddFavorite(ctx, postID, userEmail)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}

func (ps *PostService) GetAllFavourites(ctx context.Context, userEmail string, limit, offset uint) (*[]dto.PostDto, *exceptions.Error_) {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return nil, exc
	}
	if !*userExists {
		return nil, &exceptions.ErrUserNotFound
	}
	posts, err := ps.pr.GetAllFavourite(ctx, userEmail, limit, offset)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return posts, nil
}

func (ps *PostService) SearchByTitleOrAuthorOrGenre(ctx context.Context, query string, limit, offset uint, userEmail string) (*[]dto.PostDto, *exceptions.Error_) {
	posts, err := ps.pr.SearchByTitleOrAuthorOrGenre(ctx, query, limit, offset, userEmail)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return posts, nil
}

func (ps *PostService) GetAllMyBooked(ctx context.Context, userEmail string, limit, offset uint) (*[]dto.PostDto, *exceptions.Error_) {
	userExists, exc := ps.us.IsUserExist(ctx, userEmail, "")
	if exc != nil {
		return nil, exc
	}
	if !*userExists {
		return nil, &exceptions.ErrUserNotFound
	}
	posts, err := ps.pr.GetAllMyBooked(ctx, userEmail, limit, offset)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return posts, nil
}

func (s *PostService) getImageType(image *multipart.FileHeader) string {
	imageType := strings.Split(image.Filename, ".")
	return imageType[len(imageType)-1]
}

func (s *PostService) AddImage(
	ctx context.Context, postID int64, header *multipart.FileHeader, image multipart.File, userEmail string,
) (string, *exceptions.Error_) {
	post, err := s.pr.Get(ctx, postID, userEmail)
	if err != nil {
		return "", &exceptions.ErrDatabaseError
	}
	if post == nil {
		return "", &exceptions.PostNotFoudErr
	}

	imageType := s.getImageType(header)
	if imageType != "jpg" && imageType != "jpeg" && imageType != "png" {
		return "", &exceptions.ErrUnsupportedImageType
	}
	uri, err := s.pr.AddImage(ctx, image, header.Filename, imageType)
	if err != nil {
		return "", &exceptions.ErrDatabaseError
	}
	err = s.pr.UpdateImageURL(ctx, postID, uri)
	if err != nil {
		return "", &exceptions.ErrDatabaseError
	}
	return uri, nil
}

func (s *PostService) GenerateBriefContent(ctx context.Context, postID int64, userEmail string) (string, *exceptions.Error_) {
	post, err := s.GetPost(ctx, postID, userEmail)
	if err != nil {
		return "", err
	}
	content, err := s.gs.GenerateBriefContent(ctx, post.Title, post.Author)
	if err != nil {
		return "", &exceptions.ErrDatabaseError
	}
	return content, nil
}

func (s *PostService) GenerateQuote(ctx context.Context, postID int64, userEmail string) (string, *exceptions.Error_) {
	post, err := s.GetPost(ctx, postID, userEmail)
	if err != nil {
		return "", err
	}
	content, err := s.gs.GenerateQuote(ctx, post.Title, post.Author)
	if err != nil {
		return "", &exceptions.ErrDatabaseError
	}
	return content, nil
}
