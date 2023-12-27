package delivery

import (
	"context"
	"nox-ai/domain/usecase"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	ctx context.Context
	usecase usecase.IUsecase
}

func NewDelivery(e *echo.Echo, usecase usecase.IUsecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}