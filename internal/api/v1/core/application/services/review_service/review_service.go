package review_service

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/core/application/services/user_service"
)

type ReviewService struct {
	rr repositories.ReviewRepository
	us user_service.UserService
}

func NewReviewService(rr *repositories.ReviewRepository, us *user_service.UserService) *ReviewService {
	return &ReviewService{rr: *rr, us: *us}
}

func (rs *ReviewService) CreateReview(ctx context.Context, reviewerEmail string, review *dto.ReviewToCreateDto) *exceptions.Error_ {
	if userExists, exc := rs.us.IsUserExist(ctx, review.TargetUserEmail, ""); !*userExists {
		return &exceptions.ErrUserNotFound
	} else if exc != nil {
		return exc
	}

	if userExists, exc := rs.us.IsUserExist(ctx, reviewerEmail, ""); !*userExists {
		return &exceptions.ErrUserNotFound
	} else if exc != nil {
		return exc
	}

	if review.TargetUserEmail == reviewerEmail {
		return &exceptions.ErrUserIsTryingToReviewHimself
	}

	if exc := rs.rr.Create(ctx, &dto.ReviewWithoutIDDto{
		TargetUserEmail:   review.TargetUserEmail,
		ReviewerUserEmail: reviewerEmail,
		Rating:            *review.Rating,
		Comment:           *review.Comment,
		CreatedAt:         time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}); exc != nil {
		return &exceptions.ErrDatabaseError
	}

	return nil
}

func (rs *ReviewService) GetReviewsForUser(ctx context.Context, targetUsername string, limit, offset uint) (*[]dto.ReviewToGetDto, *exceptions.Error_) {
	user, exc := rs.us.GetUserByUsername(ctx, targetUsername)
	if exc != nil {
		return nil, exc
	}
	reviews, err := rs.rr.GetReviewsForUser(ctx, user.Email, limit, offset)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	return reviews, nil
}
