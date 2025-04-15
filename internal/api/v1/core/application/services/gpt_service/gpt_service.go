package gpt_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
	"example.com/m/internal/config"
	"go.uber.org/zap"
)

type GPTService struct {
	catalogID string
	logger    *zap.Logger
	cr        repositories.ChatRepository
}

func NewGPTService(catalogID string, logger *zap.Logger, cr *repositories.ChatRepository) *GPTService {
	return &GPTService{
		catalogID: catalogID,
		logger:    logger,
		cr:        *cr,
	}
}

func (s *GPTService) MakeRequestToGPTPro(prompt string, userText string) (string, error) {
	url := "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

	requestPayload := dto.GPTRequest{
		ModelUri: fmt.Sprintf("gpt://%s/yandexgpt/rc", s.catalogID),
		CompletionOptions: dto.CompletionOptions{
			Stream:      false,
			Temperature: 0.3,
			MaxTokens:   "2000",
			ReasoningOptions: dto.ReasoningOptions{
				Mode: "DISABLED",
			},
		},
		Messages: []dto.Message{
			{
				Role: "system",
				Text: prompt,
			},
			{
				Role: "user",
				Text: userText,
			},
		},
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.IAMToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gptResp dto.GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return "", err
	}

	if len(gptResp.Result.Alternatives) > 0 {
		return gptResp.Result.Alternatives[0].Message.Text, nil
	}

	return "", errors.New("Нет ответа.")
}

func (s *GPTService) GenerateBriefContent(ctx context.Context, bookName string, author string) (string, *exceptions.Error_) {
	result, err := s.MakeRequestToGPTPro(
		"Перескажи книгу"+
			"кратко и по сути. Выдели основные идеи, ключевые события и главные выводы."+
			"Изложение должно быть четким, логичным и без лишних деталей."+
			"Если такой книги нет, то напиши об этом.",
		fmt.Sprintf("Книга: %s, Автор: %s", bookName, author),
	)
	if err != nil {
		return "", &exceptions.ErrServiceUnavailable
	}
	return result, nil
}

func (s *GPTService) GenerateQuote(ctx context.Context, bookName string, author string) (string, *exceptions.Error_) {
	result, err := s.MakeRequestToGPTPro(
		"Напиши 1 настоящую цитату из книги. Напиши только цитату и ничего более."+
			"Если такой книги нет, то напиши об этом.",
		fmt.Sprintf("Книга: %s, Автор: %s", bookName, author),
	)
	if err != nil {
		return "", &exceptions.ErrServiceUnavailable
	}
	return result, nil
}

func convertMessagesToDto(messages ...dto.ChatMessage) []dto.YandexMessageDto {
	res := make([]dto.YandexMessageDto, 0, len(messages))
	for _, el := range messages {
		if el.Writer == dto.USER {
			res = append(res, dto.YandexMessageDto{
				Role: "user",
				Text: el.Message,
			})
		} else {
			res = append(res, dto.YandexMessageDto{
				Role: "assistant",
				Text: el.Message,
			})
		}
	}

	return res
}

func convertDtoToMessages(messagesDto ...dto.YandexAlternativeDto) []dto.ChatMessage {
	res := make([]dto.ChatMessage, 0, len(messagesDto))
	for _, el := range messagesDto {
		if el.Message.Role == "user" {
			res = append(res, dto.ChatMessage{
				Writer:    dto.USER,
				Message:   el.Message.Text,
				CreatedAt: time.Now().Format("02.01.2006"),
			})
		} else if el.Message.Role == "assistant" {
			res = append(res, dto.ChatMessage{
				Writer:    dto.BOT,
				Message:   el.Message.Text,
				CreatedAt: time.Now().Format("02.01.2006"),
			})
		}
	}

	return res
}

func (r *GPTService) Chat(ctx context.Context, messages ...dto.ChatMessage) ([]dto.ChatMessage, error) {
	messagesDto := make([]dto.YandexMessageDto, 0, 2)
	messagesDto = append(messagesDto, dto.YandexMessageDto{
		Role: "system",
		Text: config.Config.ChatBotPrompt,
	})
	messagesDto = append(messagesDto, convertMessagesToDto(messages...)...)
	optDto := dto.YandexCompletionOptionsDto{
		MaxTokens:   2000,
		Temperature: 0.3,
	}
	reqDto := dto.YandexRequstDto{
		ModelUri:          fmt.Sprintf("gpt://%s/yandexgpt/rc", r.catalogID),
		CompletionOptions: optDto,
		Messages:          messagesDto,
	}

	jsonData, err := json.Marshal(reqDto)
	if err != nil {
		return []dto.ChatMessage{}, err
	}

	url := "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return []dto.ChatMessage{}, err
	}

	req.Header.Set("Authorization", "Bearer "+config.IAMToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []dto.ChatMessage{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []dto.ChatMessage{}, err
	}

	var responseDTO dto.YandexResponseDto
	err = json.Unmarshal(body, &responseDTO)
	if err != nil {
		return []dto.ChatMessage{}, err
	}
	r.logger.Error(
		"Bot Repository Error",
		zap.String("method", "Chat"),
		zap.Any("resp", responseDTO),
	)

	return convertDtoToMessages(responseDTO.Result.Alternatives...), nil

}

func (s *GPTService) GetAnswer(ctx context.Context, email string, message *dto.ChatMessage) (*dto.ChatMessage, *exceptions.Error_) {
	if message == nil || message.Writer != dto.USER {
		return nil, &exceptions.ErrInvalidQuestion
	}
	message.CreatedAt = time.Now().Format("02.01.2006")

	fmt.Println(1)
	lastMessages, err := s.cr.GetChatByUserEmail(ctx, email, 0, 20)
	if err != nil {
		return nil, &exceptions.InternalServerError
	}
	fmt.Println(2)
	lastMessages = append(lastMessages, *message)
	if err := s.cr.AddNewMessageToChatByEmail(ctx, email, message); err != nil {
		return nil, &exceptions.InternalServerError
	}

	fmt.Println(3)
	answer, err := s.Chat(ctx, lastMessages...)
	if err != nil {
		return nil, &exceptions.InternalServerError
	}

	if len(answer) == 0 {
		return nil, &exceptions.ErrServiceUnavailable
	}

	if err := s.cr.AddNewMessageToChatByEmail(ctx, email, &answer[len(answer)-1]); err != nil {
		return nil, &exceptions.InternalServerError
	}

	return &answer[len(answer)-1], nil
}

func (s *GPTService) GetChatByEmail(ctx context.Context, email string, offset int, limit int) ([]dto.ChatMessage, *exceptions.Error_) {
	res, err := s.cr.GetChatByUserEmail(ctx, email, offset, limit)
	if err != nil {
		return []dto.ChatMessage{}, &exceptions.InternalServerError
	}

	return res, nil
}
