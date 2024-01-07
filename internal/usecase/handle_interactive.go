package usecase

import (
	"context"
	"nox-ai/domain/entity"
	"nox-ai/domain/service"
	"nox-ai/internal/service/model"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func (u *Usecase) HandleInteractive(ctx context.Context, user *entity.User, messageId, replyId string) error {
	err := u.service.MarkRead(ctx, model.WhatsAppStatus{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageId,
	})
	if err != nil {
		u.logger.Error("Error mark read message", zap.Error(err))
		panic(err)
	}

	switch replyId {
	case "new-chat":
		if err := NewChat(ctx, u.redis, u.service, user); err != nil {
			u.logger.Error("Error send message", zap.Error(err))
			return err
		}
	case "my-account":
		if err := MyAccount(ctx, u.service, user); err != nil {
			u.logger.Error("Error send message", zap.Error(err))
			return err
		}
	}

	return nil
}

func NewChat(ctx context.Context, redis *redis.Client, service service.IService, user *entity.User) error {
	err := deleteRedis(ctx, redis, user.Number+":prompt")
	if err != nil {
		return err
	}

	err = service.SendWA(ctx, user.Number, model.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               user.Number,
		Type:             "text",
		Text: model.MessageText{
			PreviewURL: false,
			Body:       "Your chat history with system has been reset",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func MyAccount(ctx context.Context, service service.IService, user *entity.User) error {
	var message string
	switch user.Plan {
	case entity.Free:
		message = "Your account is Free plan. You are using GPT 3.5"
	case entity.Basic:
		message = "Your account is Basic plan. You are using GPT 4 with limit 100 prompt per day"
	case entity.Premium:
		message = "Your account is Premium plan. You are using GPT 4 with no limit"
	default:
		message = "Your account is not registered"
	}

	err := service.SendWA(ctx, user.Number, model.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               user.Number,
		Type:             "text",
		Text: model.MessageText{
			PreviewURL: false,
			Body:       message,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteRedis(ctx context.Context, redis *redis.Client, key string) error {
	err := redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
