package repositories

import (
	"context"

	"example.com/m/internal/api/v1/core/application/dto"
	"firebase.google.com/go/v4/messaging"
	"github.com/appleboy/go-fcm"
	"go.uber.org/zap"
)

type FcmRepository struct {
	fc     *fcm.Client
	logger *zap.Logger
}

func NewFcmRepository(fc *fcm.Client, logger *zap.Logger) *FcmRepository {
	return &FcmRepository{
		fc:     fc,
		logger: logger,
	}
}

func notificationDtoToMessageData(notification *dto.NotificationDto) map[string]string {
	m := map[string]string{
		"title":  notification.Title,
		"conent": notification.Content,
		"image":  notification.Image,
	}

	return m
}

func (fr *FcmRepository) SendByToken(ctx context.Context, token string, notification *dto.NotificationDto) error {
	resp, err := fr.fc.Send(
		ctx,
		&messaging.Message{
			Token: token,
			Data:  notificationDtoToMessageData(notification),
		},
	)
	fr.logger.Info(
		"Fcm Repository Info",
		zap.String("method", "SendByToken"),
		zap.Any("resp", resp),
	)
	if err != nil {
		fr.logger.Error(
			"Fcm Repository Error",
			zap.String("method", "SendByToken"),
			zap.String("error", err.Error()),
		)

		return err
	}

	return nil
}
