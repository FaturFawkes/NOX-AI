package delivery

import (
	"nox-ai/domain/usecase"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Delivery struct {
	usecase usecase.IUsecase
	logger  *zap.Logger
}

func NewDelivery(e *echo.Echo, usecase usecase.IUsecase, logger *zap.Logger) *Delivery {
	return &Delivery{
		usecase: usecase,
		logger:  logger,
	}
}
