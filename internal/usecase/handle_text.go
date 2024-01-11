package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"nox-ai/domain/entity"
	"nox-ai/internal/service/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (u *Usecase) HandleText(ctx context.Context, user *entity.User, messageId, text string) error {
	var prompt []openai.ChatCompletionMessage

	err := u.service.MarkRead(ctx, model.WhatsAppStatus{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageId,
	})
	if err != nil {
		u.logger.Error("Error mark read message", zap.Error(err))
		panic(err)
	}

	if text == "/menu" {
		err = u.service.SendWA(ctx, model.InteractiveMessage{
			MessagingProduct: "whatsapp",
			RecipientType: "individual",
			To: user.Number,
			Type: "interactive",
			Interactive: model.InteractiveData{
				Type: "list",
				Body: model.InteractiveText{
					Text: "Silahkan Pilih Menu Berikut",
				},
				Action: model.InteractiveAction{
					Button: "Menu",
					Sections: []model.InteractiveSection{
						{
							Title: "Menu",
							Rows: []model.InteractiveRow{
								{
									ID: "new-chat",
									Title: "New Chat",
								},
								{
									ID: "my-account",
									Title: "My Account",
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			u.logger.Error("Error sending message", zap.Error(err))
			return err
		}
	} else {

		// Get history gpt user
		promptRedis, err := getRedis(ctx, u.redis, user.Number+":prompt")
		if err != nil {
			u.logger.Info("No history from redis", zap.Error(err))
		}

		if promptRedis != "" {
			err = json.Unmarshal([]byte(promptRedis), &prompt)
			if err != nil {
				u.logger.Error("Error unmarshal prompt group", zap.Error(err))
			}
		}

		// Add prompt user
		prompt = append(prompt, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		resGpt, err := u.service.TextGPT(ctx, prompt)
		if err != nil {
			u.logger.Error("Error generate gpt", zap.Error(err))
			return errors.New("error gpt")
		}

		// Add prompt system
		prompt = append(prompt, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resGpt,
		})

		// Set prompt to redis
		err = setRedis(ctx, u.redis, user.Number+":prompt", prompt, 0)
		if err != nil {
			u.logger.Error("Error set redis", zap.Error(err))
			return err
		}
		
		err = u.service.SendWA(ctx, model.WhatsAppMessage{
			MessagingProduct: "whatsapp",
			RecipientType: "individual",
			To: user.Number,
			Type: "text",
			Text: model.MessageText{
				PreviewURL: false,
				Body: resGpt,
			},
		})
		if err != nil {
			u.logger.Error("Error sending message", zap.Error(err))
			return err
		}
	}

	return nil
}

func getRedis(ctx context.Context, redis *redis.Client, key string) (string, error) {
	res, err := redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func setRedis(ctx context.Context, redis *redis.Client, key string, data any, duration time.Duration) error {

	dataCnv, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = redis.Set(ctx, key, string(dataCnv), duration).Err()
	if err != nil {
		return err
	}

	return nil
}