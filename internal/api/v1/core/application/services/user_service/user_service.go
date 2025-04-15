package user_service

import (
	"context"
	"fmt"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/utils"
)

type UserService struct {
	ur  repositories.UserRepository
	ptr *repositories.PushTokenRepository
}

func NewUserService(ur *repositories.UserRepository, ptr *repositories.PushTokenRepository) *UserService {
	return &UserService{ur: *ur, ptr: ptr}
}

func (s *UserService) BindPushToken(ctx context.Context, email string, token string) *exceptions.Error_ {
	if err := s.ptr.Set(&ctx, email, token); err != nil {
		return &exceptions.InternalServerError
	}

	return nil
}

func (s *UserService) IsUserExist(ctx context.Context, email string, username string) (*bool, *exceptions.Error_) {
	foundUserByEmail, err := s.ur.GetByEmail(ctx, &email)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	foundUserByUsername, err := s.ur.GetByUsername(ctx, &username)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	state := foundUserByEmail != nil || foundUserByUsername != nil
	return &state, nil
}

func (s *UserService) GetUserByUsernameWithRating(ctx context.Context, username string) (*dto.UserWithRatingDto, *exceptions.Error_) {
	user, err := s.ur.GetByUsername(ctx, &username)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	if user == nil {
		return nil, &exceptions.ErrUserNotFound
	}

	rating, err := s.ur.GetAverageReviewRatingByEmail(ctx, &user.Email)
	if err != nil {
		fmt.Println(err)
		return nil, &exceptions.ErrDatabaseError
	}

	return &dto.UserWithRatingDto{
		Username:         username,
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		TelegramUsername: user.TelegramUsername,
		IsAdmin:          user.IsAdmin,
		Rating:           rating,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, u dto.CreateUserDto) (*dto.UserDto, *exceptions.Error_) {
	userExists, exception := s.IsUserExist(ctx, u.Email, u.Username)
	if exception != nil {
		return nil, exception
	}
	if *userExists {
		return nil, &exceptions.ErrUserAlreadyExists
	}

	hashedPassword, _ := utils.HashPassword(u.Password)
	userToCreate := dto.UserDto{
		Email:            u.Email,
		Password:         hashedPassword,
		Username:         u.Username,
		CreatedAt:        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		TelegramUsername: u.TelegramUsername,
		IsAdmin:          false,
	}

	err := s.ur.Create(ctx, &userToCreate)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	return &userToCreate, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*dto.UserDto, *exceptions.Error_) {
	user, err := s.ur.GetByEmail(ctx, &email)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	if user == nil {
		return nil, &exceptions.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserDto, *exceptions.Error_) {
	user, err := s.ur.GetByUsername(ctx, &username)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	if user == nil {
		return nil, &exceptions.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) UpdateUserByEmail(ctx context.Context, email string, u dto.UpdateUserDto) (*dto.UserDto, *exceptions.Error_) {
	userExists, exception := s.IsUserExist(ctx, "", u.Username)
	if exception != nil {
		return nil, exception
	}
	if *userExists {
		return nil, &exceptions.ErrUserAlreadyExists
	}

	_, exception = s.GetUserByEmail(ctx, email)
	if exception != nil {
		return nil, exception
	}

	utils.UpdateUserTimestamps(&u)

	if err := s.ur.UpdateByEmail(ctx, &email, &u); err != nil {
		return nil, &exceptions.ErrDatabaseError
	}

	updatedUser, exception := s.GetUserByEmail(ctx, email)
	if exception != nil {
		return nil, exception
	}

	return updatedUser, nil
}
