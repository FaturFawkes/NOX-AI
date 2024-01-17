package usecase

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
)

type IUsecase interface {
	CheckNumber(ctx context.Context, data *entity.User) (*entity.User, error)
	HandleText(ctx context.Context, user *entity.User, messageId, text string) error
	HandleInteractive(ctx context.Context, user *entity.User, messageId, replyId string) error
	HandleAudio(ctx context.Context, user *entity.User, messageId, audioId string) error
}
