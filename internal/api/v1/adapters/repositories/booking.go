package repositories

import (
	"context"
	"database/sql"
	"errors"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type BookingRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewBookingRepository(db *sql.DB, logger *zap.Logger) *BookingRepository {
	return &BookingRepository{
		db:     db,
		logger: logger,
	}
}

func (r *BookingRepository) Create(ctx context.Context, b *dto.BookingToCreateDto) error {
	query, _, _ := goqu.Insert("bookings").Rows(b).ToSQL()
	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Booking Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *BookingRepository) GetByPostID(ctx context.Context, postID int64) (*dto.BookingDto, error) {
	var booking dto.BookingDto
	query, _, _ := goqu.From("bookings").Where(goqu.Ex{
		"post_id": postID,
	}).ToSQL()
	err := r.db.QueryRow(query).Scan(&booking.ID, &booking.UserEmail, &booking.PostID, &booking.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"Booking Repository Error",
			zap.String("method", "GetPostById"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Get(ctx context.Context, id int64) (*dto.BookingDto, error) {
	var booking dto.BookingDto
	query, _, _ := goqu.From("bookings").Where(goqu.Ex{
		"id": id,
	}).ToSQL()
	err := r.db.QueryRow(query).Scan(&booking.ID, &booking.UserEmail, &booking.PostID, &booking.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"Booking Repository Error",
			zap.String("method", "Get"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Delete(ctx context.Context, id int64) error {
	query, _, _ := goqu.Delete("bookings").Where(goqu.Ex{
		"id": id,
	}).ToSQL()
	_, err := r.db.Exec(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		r.logger.Error(
			"Booking Repository Error",
			zap.String("method", "Delete"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *BookingRepository) RemoveAll(ctx context.Context, postID int64) error {
	query, _, _ := goqu.Delete("bookings").ToSQL()
	_, err := r.db.Exec(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		r.logger.Error(
			"Booking Repository Error",
			zap.String("method", "RemoveAll"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}
