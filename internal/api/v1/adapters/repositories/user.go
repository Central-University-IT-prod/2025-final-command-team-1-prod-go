package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *dto.UserDto) error {
	query, _, err := goqu.Insert("users").Rows(*u).ToSQL()
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return err
	}
	_, err = r.db.Exec(query)

	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (r *UserRepository) GetAverageReviewRatingByEmail(ctx context.Context, email *string) (float64, error) {
	query, _, err := goqu.From("reviews").
		Select(goqu.COALESCE(goqu.AVG("rating"), 0).As("average_rating")).
		Where(goqu.Ex{"target_user_email": *email}, goqu.I("rating").Gt(0)).
		ToSQL()
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetAverageReviewRatingByEmail"),
			zap.String("error", err.Error()),
		)
		return 0, err
	}

	var averageRating float64
	err = r.db.QueryRow(query).Scan(&averageRating)
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetAverageReviewRatingByEmail"),
			zap.String("error", err.Error()),
		)
		return 0, err
	}

	return averageRating, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username *string) (*dto.UserDto, error) {
	query, _, err := goqu.From("users").Where(goqu.Ex{
		"username": *username,
	}).ToSQL()
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetByUsername"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	var user dto.UserDto
	err = r.db.QueryRow(query).Scan(&user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.TelegramUsername, &user.IsAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetByUsername"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email *string) (*dto.UserDto, error) {
	query, _, err := goqu.From("users").Where(goqu.Ex{
		"email": *email,
	}).ToSQL()
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetByEmail"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	var user dto.UserDto
	err = r.db.QueryRow(query).Scan(&user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.TelegramUsername, &user.IsAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "GetByEmail"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &user, nil
}

func omitEmptyFields(uMap *map[string]interface{}) {
	for k, v := range *uMap {
		if v == nil {
			delete(*uMap, k)
		}
	}
}

func (r *UserRepository) UpdateByEmail(ctx context.Context, email *string, u *dto.UpdateUserDto) error {
	var uMap map[string]interface{}
	inrec, err := json.Marshal(*u)
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "UpdateByEmail"),
			zap.String("error", err.Error()),
		)
		return err
	}
	json.Unmarshal(inrec, &uMap)
	// omitEmptyFields(&uMap)
	var rec goqu.Record = uMap
	query, _, err := goqu.From("users").Where(goqu.C("email").Eq(*email)).Update().Set(
		rec,
	).ToSQL()
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "UpdateByEmail"),
			zap.String("error", err.Error()),
		)
		return err
	}

	_, err = r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"User Repository Error",
			zap.String("method", "UpdateByEmail"),
			zap.String("error", err.Error()),
		)
		return err
	}

	return nil
}
