package repositories

import (
	"context"
	"database/sql"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type ReviewRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewReviewRepository(db *sql.DB, logger *zap.Logger) *ReviewRepository {
	return &ReviewRepository{db: db, logger: logger}
}

func (r *ReviewRepository) Create(ctx context.Context, review *dto.ReviewWithoutIDDto) error {
	query, _, err := goqu.Insert("reviews").Rows(*review).ToSQL()
	if err != nil {
		r.logger.Error(
			"Review Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Review Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (r *ReviewRepository) GetReviewsForUser(ctx context.Context, targetUserEmail string, limit, offset uint) (*[]dto.ReviewToGetDto, error) {
	var reviews []dto.ReviewToGetDto
	query, _, err := goqu.
		Select(
			"reviews.*",
			goqu.I("users.username").As("reviewer_username"),
		).
		From("reviews").
		Join(
			goqu.T("users"),
			goqu.On(
				goqu.I("reviews.reviewer_user_email").Eq(goqu.I("users.email")),
			),
		).
		Where(goqu.Ex{
			"reviews.target_user_email": targetUserEmail,
		}).
		Limit(limit).
		Offset(offset).
		ToSQL()
	if err != nil {
		r.logger.Error(
			"Review Repository Error",
			zap.String("method", "GetReviewsForUser"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error(
			"Review Repository Error",
			zap.String("method", "GetReviewsForUser"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var review dto.ReviewToGetDto
		if err := rows.Scan(
			&review.ID, &review.TargetUserEmail, &review.ReviewerUserEmail, &review.Rating,
			&review.Comment, &review.CreatedAt, &review.ReviewerUsername); err != nil {
			r.logger.Error(
				"Review Repository Error",
				zap.String("method", "GetReviewsForUser"),
				zap.String("error", err.Error()),
			)
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Review Repository Error",
			zap.String("method", "GetReviewsForUser"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &reviews, nil
}
