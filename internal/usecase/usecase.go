package usecase

import (
	"nox-ai/domain/repository"
	"nox-ai/domain/usecase"

	"go.uber.org/zap"
)

type Usecase struct {
	repo repository.IRepository
	logger *zap.Logger
}

func NewUsecase(repo repository.IRepository, logger *zap.Logger) usecase.IUsecase {
	return &Usecase{
		repo: repo,
		logger: logger,
	}
}