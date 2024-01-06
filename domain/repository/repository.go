package repository

import "nox-ai/domain/entity"

type IRepository interface {
	GetUser(number string) (*entity.User, error)
	InsertUser(data *entity.User) (*entity.User, error)
}