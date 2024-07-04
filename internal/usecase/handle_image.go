package usecase

import (
	"context"
	"encoding/base64"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/internal/delivery/request"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

func (u *Usecase) HandleImage(ctx context.Context, user *entity.User, image request.Image) error {
	if user.Plan == entity.Free || time.Now().After(*user.ExpiredAt) {
		return u.service.SendWA(model.WhatsAppMessage{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               user.Number,
			Type:             "text",
			Text: model.MessageText{
				PreviewURL: false,
				Body:       "Your are in free plan. Please upgrade to a starter or premium plan to access this feature.",
			},
		})
	}

	imageUrl, err := u.service.RetrieveMedia(image.ID)
	if err != nil {
		u.logger.Error("Error get image from whatsapp", zap.Error(err))
		return err
	}

	imagePath, err := u.service.DownloadMedia(imageUrl, entity.TypeImage)
	if err != nil {
		u.logger.Error("Error download audio", zap.Error(err))
		return err
	}

	// Getting the base64 string
	base64Image, err := encodeImage(imagePath)
	if err != nil {
		u.logger.Error("Error encoding image", zap.Error(err))
		return err
	}

	resGpt, err := u.service.TextGPT(ctx, openai.GPT4o20240513, []openai.ChatCompletionMessage{
		{
			Role: "user",
			MultiContent: []openai.ChatMessagePart{
				{
					Type: "text",
					Text: image.Caption,
				},
				{
					Type: "image_url",
					ImageURL: &openai.ChatMessageImageURL{
						URL: "data:image/png;base64," + base64Image,
					},
				},
			},
			Name:         "",
			FunctionCall: nil,
			ToolCalls:    nil,
			ToolCallID:   "",
		},
	})
	if err != nil {
		u.logger.Error("Error GPT Vision", zap.Error(err))
		return err
	}

	// Delete downloaded audio file
	err = os.Remove(imagePath)
	if err != nil {
		u.logger.Error("Error deleting image file", zap.Error(err))
		return err
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

	return nil
}

// Function to encode the image
func encodeImage(imagePath string) (string, error) {
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imageBytes), nil
}

func ImageGPT(service service.IService, text string, user *entity.User, logger *zap.Logger) error {
	prompt := strings.Split(text, "/image")
	if user.Plan == entity.Premium && time.Now().Before(*user.ExpiredAt) {
		resGptImg, err := service.ImageGPT(context.Background(), prompt[1])
		if err != nil {
			logger.Error("Error generate image", zap.Error(err))
			return err
		}

		err = service.SendWA(model.ImageMessage{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               user.Number,
			Type:             "image",
			Image: model.Image{
				Link: resGptImg.Data[0].URL,
			},
		})
		if err != nil {
			logger.Error("Error sending image", zap.Error(err))
			return err
		}
	} else {
		err := service.SendWA(model.WhatsAppMessage{
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
			logger.Error("Error sending image", zap.Error(err))
			return err
		}
	}

	return nil
}
