package migration

import (
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"gorm.io/gorm"
)

func AutoMigrate(DB *gorm.DB) {
	DB.AutoMigrate(&entity.User{})
	DB.AutoMigrate(&entity.UserLog{})
}
