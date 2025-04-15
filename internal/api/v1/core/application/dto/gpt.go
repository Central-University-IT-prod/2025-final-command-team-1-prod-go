package dto

type ReasoningOptions struct {
	Mode string `json:"mode"`
}

type CompletionOptions struct {
	Stream           bool             `json:"stream"`
	Temperature      float64          `json:"temperature"`
	MaxTokens        string           `json:"maxTokens"`
	ReasoningOptions ReasoningOptions `json:"reasoningOptions"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type GPTRequest struct {
	ModelUri          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

type Usage struct {
	InputTextTokens  string `json:"inputTextTokens"`
	CompletionTokens string `json:"completionTokens"`
	TotalTokens      string `json:"totalTokens"`
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Usage        Usage         `json:"usage"`
	ModelVersion string        `json:"modelVersion"`
}

type GPTResponse struct {
	Result Result `json:"result"`
}

