package repositories

import (
	"context"
	"errors"

	"example.com/m/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TokenRepository struct {
	rdb    *redis.Client
	logger *zap.Logger
}

func NewTokenRepository(rdb *redis.Client, logger *zap.Logger) *TokenRepository {
	return &TokenRepository{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *TokenRepository) GetByEmail(ctx *context.Context, email string) (*string, error) {
	token, err := r.rdb.Get(*ctx, email).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		r.logger.Error(
			"Token Repository Error",
			zap.String("method", "GetByEmail"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) Set(ctx *context.Context, email string, token string) error {
	err := r.rdb.Set(*ctx, email, token, config.Config.JWTExpiration).Err()
	if err != nil {
		r.logger.Error(
			"Token Repository Error",
			zap.String("method", "Set"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *TokenRepository) DeleteByEmail(ctx *context.Context, email string) error {
	err := r.rdb.Del(*ctx, email).Err()
	if err != nil {
		r.logger.Error(
			"Token Repository Error",
			zap.String("method", "Delete"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}
