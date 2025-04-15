package booking_service

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/api/v1/core/application/services/post_service"
	"example.com/m/internal/api/v1/core/application/services/user_service"
)

type BookingService struct {
	br  repositories.BookingRepository
	us  user_service.UserService
	ps  post_service.PostService
	fr  *repositories.FcmRepository
	ptr *repositories.PushTokenRepository
}

func NewBookingService(br repositories.BookingRepository, us user_service.UserService, ps post_service.PostService, fr *repositories.FcmRepository, ptr *repositories.PushTokenRepository) *BookingService {
	return &BookingService{
		br:  br,
		us:  us,
		ps:  ps,
		fr:  fr,
		ptr: ptr,
	}
}

func (bs *BookingService) BookBook(ctx context.Context, userEmail string, postID int64) *exceptions.Error_ {
	// Check if user exists
	_, exc := bs.us.GetUserByEmail(ctx, userEmail)
	if exc != nil {
		return exc
	}

	// Check if post exists
	post, exc := bs.ps.GetPost(ctx, postID, userEmail)
	if exc != nil {
		return exc
	}

	if post.Status != "available" {
		return &exceptions.ErrPostIsNotAvailable
	}

	if post.UserEmail == userEmail {
		return &exceptions.ErrUserIsOwner
	}

	// Check if this post has been booked
	existingBooking, err := bs.br.GetByPostID(ctx, postID)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	if existingBooking != nil {
		return &exceptions.ErrBookingAlreadyExists
	}

	// Create a booking
	booking := &dto.BookingToCreateDto{
		UserEmail: userEmail,
		PostID:    postID,
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}

	err = bs.br.Create(ctx, booking)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	// changing post status to booked
	_, exc = bs.ps.UpdatePost(ctx, post.UserEmail, postID, &dto.UpdatePostDto{
		Status: "booked",
	})
	if exc != nil {
		return exc
	}

	token, err := bs.ptr.GetByEmail(&ctx, post.UserEmail)
	if err != nil || token == nil {
		return nil
	}

	bs.fr.SendByToken(
		ctx,
		*token,
		&dto.NotificationDto{
			Title:   "Новая бронь!",
			Content: "Книгу \"" + post.Title + "\" забронировали.",
		},
	)

	return nil
}

func (bs *BookingService) DeleteBooking(ctx context.Context, userEmail string, postID int64) *exceptions.Error_ {
	// Check if user has already booked this post
	existingBooking, err := bs.br.GetByPostID(ctx, postID)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	if existingBooking == nil {
		return &exceptions.ErrBookingNotFound
	}

	// Check if user exists
	_, exc := bs.us.GetUserByEmail(ctx, userEmail)
	if exc != nil {
		return exc
	}

	post, exc := bs.ps.GetPost(ctx, postID, userEmail)
	if exc != nil {
		return exc
	}

	// changing post status to taken
	_, exc = bs.ps.UpdatePost(ctx, post.UserEmail, postID, &dto.UpdatePostDto{
		Status: "available",
	})
	if exc != nil {
		return exc
	}

	// Delete the booking
	err = bs.br.Delete(ctx, existingBooking.ID)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	return nil
}

func (bs *BookingService) MarkAsTaken(ctx context.Context, userEmail string, postID int64) *exceptions.Error_ {
	// Check if user exists
	_, exc := bs.us.GetUserByEmail(ctx, userEmail)
	if exc != nil {
		return exc
	}

	// Check if post exists
	post, exc := bs.ps.GetPost(ctx, postID, userEmail)
	if exc != nil {
		return exc
	}

	if post.UserEmail == userEmail {
		return &exceptions.ErrUserIsOwner
	}

	// changing post status to taken
	_, exc = bs.ps.UpdatePost(ctx, post.UserEmail, postID, &dto.UpdatePostDto{
		Status: "taken",
	})
	if exc != nil {
		return exc
	}

	// removing all bookings for this post
	err := bs.br.RemoveAll(ctx, postID)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}

	return nil
}
