package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (u *Usecase) HandleText(ctx context.Context, user *entity.User, messageId, text string) error {
	var prompt []openai.ChatCompletionMessage

	if text == "/menu" {
		return sendMenu(u.service, user, u.logger)
	} else if strings.Contains(text, "/image") {
		return ImageGPT(u.service, text, user, u.logger)
	} else {
		if user.Plan == entity.Free {
			if user.RemainingRequest == 0 {
				err := u.service.SendWA(model.WhatsAppMessage{
					MessagingProduct: "whatsapp",
					RecipientType:    "individual",
					To:               user.Number,
					Type:             "text",
					Text: model.MessageText{
						PreviewURL: false,
						Body:       "Your reach the daily limit for free tier",
					},
				})
				if err != nil {
					u.logger.Error("Error send message limit", zap.Error(err))
					return err
				}
				return nil
			} else {
				user.RemainingRequest -= 1
				err := u.repo.UpdateUser(user)
				if err != nil {
					u.logger.Error("Error decrease remaining request", zap.Error(err))
					return err
				}
			}
		}

		summarize := strings.Contains(text, "youtube.com") || strings.Contains(text, "youtu.be")
		var gptVersion string
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
		} else {
			if user.Plan != entity.Free {
				prompt = append(prompt, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: "gue adalah model gpt bisa pakai bahasa santai dan asik buat ngobrol",
				})
			}
		}

		if summarize {
			if user.Plan != entity.Free && time.Now().Before(*user.ExpiredAt) {
				code, transcribe, err := u.service.TranscribeYoutube(text, "id")
				if code != 200 {
					code, transcribe, err = u.service.TranscribeYoutube(text, "en")
				}
				if err != nil {
					u.logger.Error("Error transcribe video", zap.Error(err))
					return err
				}
				text = fmt.Sprint("Please summarize this transcript from youtube: \n", transcribe)
			} else {
				err = u.service.SendWA(model.WhatsAppMessage{
					MessagingProduct: "whatsapp",
					RecipientType:    "individual",
					To:               user.Number,
					Type:             "text",
					Text: model.MessageText{
						PreviewURL: false,
						Body:       "Your are in free plan. Please upgrade to a starter or premium plan to access this feature.",
					},
				})
				if err != nil {
					u.logger.Error("Error send message expired", zap.Error(err))
					return err
				}
				return nil
			}
		}

		// Add prompt user before gpt
		prompt = append(prompt, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		if user.Plan != entity.Free && time.Now().Before(*user.ExpiredAt) {
			gptVersion = openai.GPT4TurboPreview
		} else {
			gptVersion = openai.GPT3Dot5Turbo1106
		}

		resGpt, err := u.service.TextGPT(ctx, gptVersion, prompt)
		if err != nil {
			u.logger.Error("Error generate gpt", zap.Error(err))
			return errors.New("error gpt")
		}

		err = u.service.SendWA(model.WhatsAppMessage{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               user.Number,
			Type:             "text",
			Context: model.ContextMessage{
				MessageID: messageId,
			},
			Text: model.MessageText{
				PreviewURL: false,
				Body:       resGpt.Choices[0].Message.Content,
			},
		})
		if err != nil {
			u.logger.Error("Error sending message", zap.Error(err))
			return err
		}

		// Add prompt system after gpt
		prompt = append(prompt, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resGpt.Choices[0].Message.Content,
		})

		// Set prompt to redis if not summarize youtube
		if !summarize {
			err = setRedis(ctx, u.redis, user.Number+":prompt", prompt, 0)
			if err != nil {
				u.logger.Error("Error set redis", zap.Error(err))
				return err
			}
		}

		//	Insert Logs
		err = u.InserUserLog(user, resGpt)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendMenu(service service.IService, user *entity.User, logger *zap.Logger) error {
	err := service.SendWA(model.InteractiveMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               user.Number,
		Type:             "interactive",
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
								ID:    "new-chat",
								Title: "New Chat",
							},
							{
								ID:    "my-account",
								Title: "My Account",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		logger.Error("Error sending message", zap.Error(err))
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
