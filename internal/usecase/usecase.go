package usecase

import (
	"github.com/FaturFawkes/NOX-AI/domain/repository"
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/domain/usecase"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Usecase struct {
	repo    repository.IRepository
	logger  *zap.Logger
	service service.IService
	redis   *redis.Client
}

func NewUsecase(repo repository.IRepository, redis *redis.Client, logger *zap.Logger, service service.IService) usecase.IUsecase {
	return &Usecase{
		repo:    repo,
		service: service,
		logger:  logger,
		redis:   redis,
	}
}
