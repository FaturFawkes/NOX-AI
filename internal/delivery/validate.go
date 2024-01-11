package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (dlv *Delivery) Validate(c echo.Context) error {

	token := c.QueryParam("hub.verify_token")
	mode := c.QueryParam("hub.mode")
	challenge := c.QueryParam("hub.challenge")

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	if token == "HAPPY" && mode == "subscribe" {
		dlv.logger.Info("WEBHOOK_VERIFIED")
		return c.HTML(http.StatusOK, challenge)
	}

	return c.JSON(http.StatusForbidden, nil)
}
