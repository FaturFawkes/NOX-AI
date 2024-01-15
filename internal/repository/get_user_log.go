package repository

import "nox-ai/domain/entity"

func (r *Repository) GetUserLog(userId uint) (*entity.UserLog, error) {
	var userLog entity.UserLog
	
	if err := r.db.Where("user_id = ? AND DATE(updated_at) = CURDATE()", userId).First(&userLog).Error; err != nil {
		return nil, err
	}
	
	return &userLog, nil
}