package repository

import "github.com/FaturFawkes/NOX-AI/domain/entity"

func (r *Repository) InsertUserLog(log *entity.UserLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return err
	}

	return nil
}
