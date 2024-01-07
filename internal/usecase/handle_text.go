package usecase

import (
	"context"
	"encoding/json"
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

	// Get history gpt user
	promptRedis, err := getRedis(ctx, u.redis, user.Number+":prompt")
	if err != nil {
		u.logger.Error("Error get data redis", zap.Error(err))
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
		return err
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

	err = u.service.SendWA(ctx, user.Number, model.InteractiveMessage{
		MessagingProduct: "whatsapp",
		RecipientType: "individual",
		To: user.Number,
		Type: "interactive",
		Interactive: model.InteractiveData{
			Type: "list",
			Body: model.InteractiveText{
				Text: resGpt,
			},
			Action: model.InteractiveAction{
				Button: "Explore",
				Sections: []model.InteractiveSection{
					{
						Title: "Menu",
						Rows: []model.InteractiveRow{
							{
								ID: "new-chat",
								Title: "New Chat",
								Description: "Create new context with Chat GPT",
							},
							{
								ID: "my-account",
								Title: "My Account",
								Description: "Describe your account and plan type",
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

	return nil
}

func getRedis(ctx context.Context, redis *redis.Client, key string) (string, error) {
	res, err := redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func deleteRedis(ctx context.Context, redis *redis.Client, key string) error {
	err := redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
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
