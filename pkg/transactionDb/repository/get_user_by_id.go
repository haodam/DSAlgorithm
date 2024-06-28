package repository

import (
	"github.com/haodam/DSAlgorithm/pkg/transactionDb/entity"
	"gorm.io/gorm"
)

type userRepository struct {
}

func (r *userRepository) GetUserById(tx *gorm.DB, id int) (entity.User, error) {
	var (
		user entity.User
		err  error
	)
	err = tx.Table(user.TableName()).Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
