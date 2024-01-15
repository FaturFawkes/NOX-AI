package service

import (
	"context"
	"nox-ai/internal/service/model"

	"github.com/sashabaranov/go-openai"
)

type IService interface {
	TextGPT(ctx context.Context, model string, message []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error)
	ImageGPT(ctx context.Context, prompt string) (*openai.ImageResponse, error)
	SendWA(ctx context.Context, data any) error
	MarkRead(ctx context.Context, data model.WhatsAppStatus) error
}
