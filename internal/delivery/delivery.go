package delivery

import (
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/domain/usecase"
	"go.uber.org/zap"
)

type Delivery struct {
	usecase usecase.IUsecase
	logger  *zap.Logger
	service service.IService
}

func NewDelivery(usecase usecase.IUsecase, logger *zap.Logger, service service.IService) *Delivery {
	return &Delivery{
		usecase: usecase,
		service: service,
		logger:  logger,
	}
}
