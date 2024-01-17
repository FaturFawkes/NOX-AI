package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (u *Usecase) HandleText(ctx context.Context, user *entity.User, messageId, text string) error {
	var prompt []openai.ChatCompletionMessage

	err := u.service.MarkRead(model.WhatsAppStatus{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageId,
	})
	if err != nil {
		u.logger.Error("Error mark read message", zap.Error(err))
		panic(err)
	}

	if text == "/menu" {
		err = u.service.SendWA(model.InteractiveMessage{
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
			u.logger.Error("Error sending message", zap.Error(err))
			return err
		}
	} else if strings.Contains(text, "/image") {
		prompt := strings.Split(text, "/image")
		if user.Plan == entity.Premium && time.Now().Before(*user.ExpiredAt) {
			resGptImg, err := u.service.ImageGPT(ctx, prompt[1])
			if err != nil {
				u.logger.Error("Error generate image", zap.Error(err))
				return err
			}

			err = u.service.SendWA(model.ImageMessage{
				MessagingProduct: "whatsapp",
				RecipientType:    "individual",
				To:               user.Number,
				Type:             "image",
				Image: model.Image{
					Link: resGptImg.Data[0].URL,
				},
			})
			if err != nil {
				u.logger.Error("Error sending image", zap.Error(err))
				return err
			}
		} else {
			err = u.service.SendWA(model.WhatsAppMessage{
				MessagingProduct: "whatsapp",
				RecipientType:    "individual",
				To:               user.Number,
				Type:             "text",
				Text: model.MessageText{
					PreviewURL: false,
					Body:       "You are in free plan and not have access to this feature",
				},
			})
			if err != nil {
				u.logger.Error("Error sending image", zap.Error(err))
				return err
			}
		}
	} else {
		fmt.Println("INI USER ", *user)
		if user.Plan == entity.Free {
			if user.RemainingRequest == 0 {
				err = u.service.SendWA(model.WhatsAppMessage{
					MessagingProduct: "whatsapp",
					RecipientType:    "individual",
					To:               user.Number,
					Type:             "text",
					Text: model.MessageText{
						PreviewURL: false,
						Body:       "Your reach the limit for free tier",
					},
				})
				if err != nil {
					u.logger.Error("Error send message limit", zap.Error(err))
					return err
				}
				return nil
			} else {
				user.RemainingRequest -= 1
				err = u.repo.UpdateUser(user)
				if err != nil {
					u.logger.Error("Error decrease remaining request", zap.Error(err))
					return err
				}
			}
		}

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
			prompt = append(prompt, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: "saya adalah model dengan gpt 4 yang memiliki pengetahuan terbaru",
			})
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

		// Set prompt to redis
		err = setRedis(ctx, u.redis, user.Number+":prompt", prompt, 0)
		if err != nil {
			u.logger.Error("Error set redis", zap.Error(err))
			return err
		}

		// Add token logger
		res, err := u.repo.GetUserLog(user.ID)
		if err != nil {
			u.logger.Info("No user log for " + user.Number)
			err = u.repo.InsertUserLog(&entity.UserLog{
				UserID:        user.ID,
				TokenRequest:  resGpt.Usage.PromptTokens,
				TokenResponse: resGpt.Usage.CompletionTokens,
				TokenUsage:    resGpt.Usage.TotalTokens,
			})
			if err != nil {
				u.logger.Error("Error insert user log", zap.Error(err))
			}
			return err
		} else {
			res.TokenRequest += resGpt.Usage.PromptTokens
			res.TokenResponse += resGpt.Usage.CompletionTokens
			res.TokenUsage += resGpt.Usage.TotalTokens
			err = u.repo.UpdateUserLog(res)
			if err != nil {
				u.logger.Error("Error update log user log", zap.Error(err))
				return err
			}
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
