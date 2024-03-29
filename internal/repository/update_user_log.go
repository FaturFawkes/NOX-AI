package repository

import "github.com/FaturFawkes/NOX-AI/domain/entity"

func (r *Repository) UpdateUserLog(log *entity.UserLog) error {
	if err := r.db.Where("user_id = ?", log.UserID).Save(log).Error; err != nil {
		return err
	}

	return nil
}
