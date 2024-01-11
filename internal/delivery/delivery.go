package delivery

import (
	"go.uber.org/zap"
	"nox-ai/domain/usecase"
)

type Delivery struct {
	usecase usecase.IUsecase
	logger  *zap.Logger
}

func NewDelivery(usecase usecase.IUsecase, logger *zap.Logger) *Delivery {
	return &Delivery{
		usecase: usecase,
		logger:  logger,
	}
}
