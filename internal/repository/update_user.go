package repository

import "github.com/FaturFawkes/NOX-AI/domain/entity"

func (r *Repository) UpdateUser(user *entity.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}
