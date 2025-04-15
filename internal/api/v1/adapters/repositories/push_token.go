package repositories

import (
	"context"
	"errors"

	"example.com/m/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type PushTokenRepository struct {
	rdb    *redis.Client
	logger *zap.Logger
}

func NewPushTokenRepository(rdb *redis.Client, logger *zap.Logger) *PushTokenRepository {
	return &PushTokenRepository{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *PushTokenRepository) Set(ctx *context.Context, email string, token string) error {
	err := r.rdb.Set(*ctx, email+"_push", token, config.Config.JWTExpiration).Err()
	r.logger.Info(
		"Psuh Token Repository",
		zap.String("method", "Set"),
		zap.Any("token", token),
	)

	if err != nil {
		r.logger.Error(
			"Push Token Repository Error",
			zap.String("method", "Set"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PushTokenRepository) GetByEmail(ctx *context.Context, email string) (*string, error) {
	token, err := r.rdb.Get(*ctx, email+"_push").Result()

	if errors.Is(err, redis.Nil) {
		r.logger.Info(
			"Psuh Token Repository",
			zap.String("method", "GetByEmail empty"),
			zap.Any("token", token),
			zap.String("email", email),
		)
		return nil, nil
	} else if err != nil {
		r.logger.Error(
			"Psuh Token Repository Error",
			zap.String("method", "GetByEmail err"),
			zap.String("error", err.Error()),
			zap.String("email", email),
		)
		return nil, err
	}
	r.logger.Info(
		"Psuh Token Repository",
		zap.String("method", "GetByEmail find"),
		zap.Any("token", token),
		zap.String("email", email),
	)

	return &token, nil
}

func (r *PushTokenRepository) DeleteByEmail(ctx *context.Context, email string) error {
	err := r.rdb.Del(*ctx, email+"_push").Err()
	if err != nil {
		r.logger.Error(
			"Psuh Token Repository Error",
			zap.String("method", "Delete"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}
