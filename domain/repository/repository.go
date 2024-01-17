package repository

import (
	"github.com/FaturFawkes/NOX-AI/domain/entity"
)

type IRepository interface {
	GetUser(number string) (*entity.User, error)
	InsertUser(data *entity.User) (*entity.User, error)
	GetUserLog(userId uint) (*entity.UserLog, error)
	InsertUserLog(log *entity.UserLog) error
	UpdateUserLog(log *entity.UserLog) error
	UpdateUser(user *entity.User) error
}
