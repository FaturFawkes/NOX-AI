package service

import (
	"context"
	"nox-ai/internal/service/model"

	"github.com/sashabaranov/go-openai"
)

type IService interface {
	TextGPT(ctx context.Context, message []openai.ChatCompletionMessage) (string, error)
	SendWA(ctx context.Context, number string, data model.WhatsAppMessage) error
	MarkRead(ctx context.Context, data model.WhatsAppStatus) error
}
