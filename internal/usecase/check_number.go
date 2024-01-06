package usecase

import (
	"context"
	"nox-ai/domain/entity"

	"go.uber.org/zap"
)

func (u *Usecase) CheckNumber(ctx context.Context, data *entity.User) (*entity.User, error) {
	user, err := u.repo.GetUser(data.Number)
	if err != nil {
		user, err = u.repo.InsertUser(data)
		if err != nil {
			u.logger.Error("Error insert user to db", zap.Error(err))
			return nil, err
		}
	}

	return user, nil
}
