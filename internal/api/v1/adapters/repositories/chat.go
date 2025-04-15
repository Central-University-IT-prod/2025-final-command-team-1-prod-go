package repositories

import (
	"context"
	"database/sql"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type ChatRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewChatRepository(db *sql.DB, logger *zap.Logger) *ChatRepository {
	return &ChatRepository{
		db:     db,
		logger: logger,
	}
}

func dbMessageToChatMessage(messages ...dto.DbMessageDto) []dto.ChatMessage {
	msgs := make([]dto.ChatMessage, 0, len(messages))
	for _, el := range messages {
		if el.Role == "user" {
			msgs = append(msgs, dto.ChatMessage{
				Message:   el.Text,
				Writer:    dto.USER,
				CreatedAt: el.CreatedAt,
			})
		} else {
			msgs = append(msgs, dto.ChatMessage{
				Message:   el.Text,
				Writer:    dto.BOT,
				CreatedAt: el.CreatedAt,
			})
		}
	}

	return msgs
}

func chatMessageToDto(message *dto.ChatMessage, email string) *dto.DbMessageDto {
	if message.Writer == dto.USER {
		return &dto.DbMessageDto{
			Email:     email,
			Text:      message.Message,
			Role:      "user",
			CreatedAt: message.CreatedAt,
		}
	} else {
		return &dto.DbMessageDto{
			Email:     email,
			Text:      message.Message,
			Role:      "bot",
			CreatedAt: message.CreatedAt,
		}
	}
}

func (r *ChatRepository) GetChatByUserEmail(ctx context.Context, email string, offset int, limit int) ([]dto.ChatMessage, error) {
	query, _, _ := goqu.From("chat_messages").
		Where(goqu.Ex{
			"email": email,
		}).
		Order(goqu.C("created_at").Desc()).
		Offset(uint(offset)).
		Limit(uint(limit)).
		ToSQL()

	r.logger.Debug(
		"Chat Repository Debug",
		zap.String("method", "GetChatByUserEmail"),
		zap.String("SQL Query", query),
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error(
			"Chat Repository Error",
			zap.String("method", "GetChatByUserEmail"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	var messages []dto.DbMessageDto
	for rows.Next() {
		var msg dto.DbMessageDto
		if err := rows.Scan(&msg.Email, &msg.Text, &msg.Role, &msg.CreatedAt); err != nil {
			r.logger.Error(
				"Chat Repository Error",
				zap.String("method", "GetChatByUserEmail"),
				zap.String("error", err.Error()),
			)
			return nil, err
		}
		r.logger.Debug(
			"Chat Repository Debug",
			zap.String("method", "GetChatByUserEmail"),
			zap.Any("Get dto from db", msg),
		)
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Chat Repository Error",
			zap.String("method", "GetChatByUserEmail"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	r.logger.Debug(
		"Chat Repository Debug",
		zap.String("method", "GetChatByUserEmail"),
		zap.Any("Convert slice msgs", messages),
	)

	return dbMessageToChatMessage(messages...), nil
}

func (r *ChatRepository) AddNewMessageToChatByEmail(ctx context.Context, email string, message *dto.ChatMessage) error {
	msgDto := chatMessageToDto(message, email)
	query, _, _ := goqu.Insert("chat_messages").Rows(*msgDto).ToSQL()
	_, err := r.db.Exec(query)

	if err != nil {
		r.logger.Error(
			"Chat Repository Error",
			zap.String("method", "AddNewMessageToChatByEmail"),
			zap.String("error", err.Error()),
		)
		return err
	}

	return nil
}
