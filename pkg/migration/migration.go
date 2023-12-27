package migration

import (
	"nox-ai/domain/entity"

	"gorm.io/gorm"
)

func AutoMigrate(DB *gorm.DB) {
	DB.AutoMigrate(&entity.User{})
}
