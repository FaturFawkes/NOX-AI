package usecase

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/delivery/request"
	"github.com/sashabaranov/go-openai"
)

type IUsecase interface {
	CheckNumber(ctx context.Context, data *entity.User) (*entity.User, error)
	InserUserLog(user *entity.User, resGpt *openai.ChatCompletionResponse) error
	HandleText(ctx context.Context, user *entity.User, messageId, text string) error
	HandleInteractive(ctx context.Context, user *entity.User, replyId string) error
	HandleAudio(ctx context.Context, user *entity.User, messageId, audioId string) error
	HandleImage(ctx context.Context, user *entity.User, image request.Image) error
}
