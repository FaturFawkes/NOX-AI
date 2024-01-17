package service

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"

	"github.com/sashabaranov/go-openai"
)

type IService interface {
	TextGPT(ctx context.Context, model string, message []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error)
	ImageGPT(ctx context.Context, prompt string) (*openai.ImageResponse, error)
	SendWA(data any) error
	MarkRead(data model.WhatsAppStatus) error
}
