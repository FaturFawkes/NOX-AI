package delivery

import (
	"context"
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/delivery/request"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (dlv *Delivery) Message(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occured", r)
			return
		}
		fmt.Println("Application running perfectly")
	}()

	var message request.WhatsAppBusinessAccount

	err := c.Bind(&message)
	if err != nil {
		dlv.logger.Error("Error binding message", zap.Error(err))
		return c.JSON(http.StatusOK, "bad request")
	}

	// Logger
	dlv.logger.Info("Received message", zap.Any("message", message))

	if message.Entry[0].Changes[0].Value.Messages == nil {
		return c.JSON(http.StatusOK, "Callback OK")
	}

	data := &entity.User{
		Name:             message.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
		Number:           message.Entry[0].Changes[0].Value.Contacts[0].WaID,
		ExpiredAt:        nil,
		Plan:             entity.Free,
		RemainingRequest: 20,
	}

	err = dlv.service.MarkRead(model.WhatsAppStatus{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        message.Entry[0].Changes[0].Value.Messages[0].ID,
	})
	if err != nil {
		fmt.Println("Error read message", err.Error())
	}

	user, err := dlv.usecase.CheckNumber(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusOK, zap.Error(err))
	}

	switch message.Entry[0].Changes[0].Value.Messages[0].Type {
	case "text":
		go func() {
			err = dlv.usecase.HandleText(context.Background(), user, message.Entry[0].Changes[0].Value.Messages[0].Text.Body)
		}()
		return c.JSON(http.StatusOK, nil)
	case "audio":
		go func() {
			err = dlv.usecase.HandleAudio(context.Background(), user, message.Entry[0].Changes[0].Value.Messages[0].ID, message.Entry[0].Changes[0].Value.Messages[0].Audio.ID)
		}()
	case "reaction":

	case "image":
		go func() {
			err = dlv.usecase.HandleImage(context.Background(), user, message.Entry[0].Changes[0].Value.Messages[0].Image)
		}()

	case "sticker":

	case "location":

	case "button":

	case "interactive":
		err = dlv.usecase.HandleInteractive(c.Request().Context(), user, message.Entry[0].Changes[0].Value.Messages[0].Interactive.ListReply.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error handle interactive message")
		}
	}
	return c.JSON(http.StatusOK, nil)
}
