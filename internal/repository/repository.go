package repository

import (
	"nox-ai/domain/repository"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository.IRepository {
	return &Repository{
		db: db,
	}
} 