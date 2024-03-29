package repository

import "github.com/FaturFawkes/NOX-AI/domain/entity"

func (r *Repository) GetUser(number string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("number = ?", number).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
