package delivery

import (
	"context"
	"fmt"
	"net/http"
	"nox-ai/domain/entity"
	"nox-ai/internal/delivery/request"

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

	data := &entity.User{
		Name:      message.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
		Number:    message.Entry[0].Changes[0].Value.Contacts[0].WaID,
		ExpiredAt: nil,
		Plan:      entity.Free,
	}

	user, err := dlv.usecase.CheckNumber(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusOK, zap.Error(err))
	}

	switch message.Entry[0].Changes[0].Value.Messages[0].Type {
	case "text":
		go func ()  {
			err = dlv.usecase.HandleText(context.Background(), user, message.Entry[0].Changes[0].Value.Messages[0].ID, message.Entry[0].Changes[0].Value.Messages[0].Text.Body)
		}()
		return c.JSON(http.StatusOK, nil)
	case "reaction":

	case "image":

	case "sticker":

	case "location":

	case "button":

	case "interactive":
		err = dlv.usecase.HandleInteractive(c.Request().Context(), user, message.Entry[0].Changes[0].Value.Messages[0].ID, message.Entry[0].Changes[0].Value.Messages[0].Interactive.ListReply.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error handle interactive message")
		}
	}
	return c.JSON(http.StatusOK, nil)
}
