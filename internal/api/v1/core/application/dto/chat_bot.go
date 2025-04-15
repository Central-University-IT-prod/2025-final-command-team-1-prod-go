package dto

type (
	MessageDto struct {
		Text string `json:"text" binding:"required"`
	}

	ChatMessageDto struct {
		Role      string `json:"role" binding:"required"`
		Text      string `json:"text" binding:"required"`
		CreatedAt string `json:"created_at"`
	}

	ChatDto struct {
		Messages []ChatMessageDto `json:"messages" binding:"required"`
	}

	DbMessageDto struct {
		Email     string `db:"email"`
		Text      string `db:"text"`
		Role      string `db:"role"`
		CreatedAt string `db:"created_at"`
	}

	ChatQueryDto struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"min=0"`
	}

	YandexRequstDto struct {
		ModelUri          string                     `json:"model_uri"`
		CompletionOptions YandexCompletionOptionsDto `json:"completion_options"`
		Messages          []YandexMessageDto         `json:"messages"`
	}

	YandexCompletionOptionsDto struct {
		MaxTokens   int     `json:"max_tokens"`
		Temperature float32 `json:"temperature"`
	}

	YandexMessageDto struct {
		Role string `json:"role"`
		Text string `json:"text"`
	}

	YandexResponseDto struct {
		Result YandexResultDto `json:"result"`
	}

	YandexResultDto struct {
		Alternatives []YandexAlternativeDto `json:"alternatives"`
	}

	YandexAlternativeDto struct {
		Message YandexMessageDto `json:"message"`
		Status  string           `json:"status"`
	}
)


type writerType string

const (
	USER writerType = "user"
	BOT  writerType = "bot"
)

type ChatMessage struct {
	Message   string
	Writer    writerType
	CreatedAt string
}