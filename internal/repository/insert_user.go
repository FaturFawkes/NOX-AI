package repository

import "github.com/FaturFawkes/NOX-AI/domain/entity"

func (r *Repository) InsertUser(data *entity.User) (*entity.User, error) {

	if err := r.db.Create(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
