package repository

import (
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"time"
)

func (r *Repository) GetUserLog(userId uint) (*entity.UserLog, error) {
	var userLog entity.UserLog

	today := time.Now().Format("2006-01-02")
	if err := r.db.Where("user_id = ? AND DATE(updated_at) = ?", userId, today).First(&userLog).Error; err != nil {
		return nil, err
	}

	return &userLog, nil
}
