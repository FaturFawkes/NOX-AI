package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"os"
)

func (u *Usecase) HandleAudio(ctx context.Context, user *entity.User, messageId, audioId string) error {
	var prompt []openai.ChatCompletionMessage

	// Get Audio Data
	audioUrl, err := u.service.RetrieveMedia(audioId)
	if err != nil {
		u.logger.Error("Error retrieve audio", zap.Error(err))
		return err
	}

	audioPath, err := u.service.DownloadMedia(audioUrl, entity.TypeAudio)
	if err != nil {
		u.logger.Error("Error download audio", zap.Error(err))
		return err
	}

	text, err := u.service.TranscriptionGPT(ctx, audioPath)
	if err != nil {
		u.logger.Error("Error transcription audio", zap.Error(err))
		return err
	}

	// Delete downloaded audio file
	err = os.Remove(audioPath)
	if err != nil {
		u.logger.Error("Error deleting audio file", zap.Error(err))
		return err
	}

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

	// Add prompt user before gpt
	prompt = append(prompt, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	// Get GPT Response
	resGpt, err := u.service.TextGPT(ctx, openai.GPT4o, prompt)
	if err != nil {
		u.logger.Error("Error generate gpt", zap.Error(err))
		return errors.New("error gpt")
	}

	// Add prompt user aftert gpt
	prompt = append(prompt, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: resGpt.Choices[0].Message.Content,
	})

	// Sending voice to Whatsapp
	err = u.service.SendWA(model.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               user.Number,
		Type:             "text",
		Text: model.MessageText{
			PreviewURL: false,
			Body:       resGpt.Choices[0].Message.Content,
		},
		Context: model.ContextMessage{
			MessageID: messageId,
		},
	})
	if err != nil {
		u.logger.Error("Error sending message", zap.Error(err))
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
