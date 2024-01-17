package delivery

import (
	"github.com/FaturFawkes/NOX-AI/domain/usecase"
	"go.uber.org/zap"
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
