package usecase

import (
	"context"
	"nox-ai/domain/entity"
	"nox-ai/domain/service"
	"nox-ai/internal/service/model"
)

func CreateImage(ctx context.Context, service service.IService, prompt string, user *entity.User) error {
	res, err := service.ImageGPT(ctx, prompt)
	if err != nil {
		return err
	}

	err = service.SendWA(ctx, model.ImageMessage{
		MessagingProduct: "whatsapp",
		RecipientType: "individual",
		To: user.Number,
		Type: "image",
		Image: model.Image{
			Link: res.Data[0].URL,
		},
	})
	if err != nil {
		return err
	}
	
	return nil
}