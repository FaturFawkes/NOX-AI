package delivery

import (
	"net/http"
	"nox-ai/domain/entity"
	"nox-ai/internal/delivery/request"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (dlv *Delivery) Message(c echo.Context) error {
	var message request.WhatsAppBusinessAccount
	
	err := c.Bind(&message)
	if err != nil {
		dlv.logger.Error("Error binding message", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	data := &entity.User{
		Name: message.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
		Number: message.Entry[0].Changes[0].Value.Contacts[0].WaID,
		ExpiredAt: nil,
		Plan: entity.Free,
	}

	user, err := dlv.usecase.CheckNumber(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, zap.Error(err))
	}

	switch message.Entry[0].Changes[0].Value.Messages[0].Type {
	case "text":
		err = dlv.usecase.HandleText(c.Request().Context(), user, message.Entry[0].Changes[0].Value.Messages[0].ID, message.Entry[0].Changes[0].Value.Messages[0].Text.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error handle message")
		}
	case "reaction":
	
	case "image":
	
	case "sticker":

	case "location":
	
	case "button":
	
	case "interactive":
		err = dlv.usecase.HandleInteractive(c.Request().Context(), user, message.Entry[0].Changes[0].Value.Messages[0].ID,  message.Entry[0].Changes[0].Value.Messages[0].Interactive.ListReply.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error handle interactive message")
		}
	}
	return nil
}
