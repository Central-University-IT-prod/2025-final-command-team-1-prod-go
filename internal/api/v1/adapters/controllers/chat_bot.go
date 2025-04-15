package controllers

import (
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/gpt_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type ChatBotController struct {
	chatBotService *gpt_service.GPTService
}

func NewChatBotController(chatBotService *gpt_service.GPTService) *ChatBotController {
	return &ChatBotController{
		chatBotService: chatBotService,
	}
}

func (c *ChatBotController) SendMessage(ctx *gin.Context) {
	var message dto.MessageDto
	if err := ctx.ShouldBindBodyWithJSON(&message); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	email := payload["email"].(string)

	ans, err := c.chatBotService.GetAnswer(ctx.Request.Context(), email, &dto.ChatMessage{
		Message: message.Text,
		Writer:  dto.USER,
	})
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(200, &dto.MessageDto{
		Text: ans.Message,
	})
}

func (c *ChatBotController) GetChat(ctx *gin.Context) {
	var query dto.ChatQueryDto
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	email := payload["email"].(string)

	chat, err := c.chatBotService.GetChatByEmail(ctx.Request.Context(), email, query.Offset, query.Limit)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	chatMessagesDto := make([]dto.ChatMessageDto, 0, len(chat))
	for _, el := range chat {
		if el.Writer == dto.USER {
			chatMessagesDto = append(chatMessagesDto, dto.ChatMessageDto{
				Role:      "user",
				Text:      el.Message,
				CreatedAt: el.CreatedAt,
			})
		} else {
			chatMessagesDto = append(chatMessagesDto, dto.ChatMessageDto{
				Role:      "bot",
				Text:      el.Message,
				CreatedAt: el.CreatedAt,
			})
		}
	}

	ctx.JSON(200, &dto.ChatDto{
		Messages: chatMessagesDto,
	})
}