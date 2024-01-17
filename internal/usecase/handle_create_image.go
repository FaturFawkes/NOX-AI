package usecase

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
)

func CreateImage(ctx context.Context, service service.IService, prompt string, user *entity.User) error {
	res, err := service.ImageGPT(ctx, prompt)
	if err != nil {
		return err
	}

	err = service.SendWA(model.ImageMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               user.Number,
		Type:             "image",
		Image: model.Image{
			Link: res.Data[0].URL,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
