package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) AddUser(user model.User) error {
	data := u.db.Create(&user)

	if data.Error != nil {
		return data.Error
	}

	return nil
}

func (u *UserRepository) UserAvail(cred model.User) error {
	result := u.db.Where("username = ?", cred.Username).Where("password = ?", cred.Password).First(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *UserRepository) CheckPassLength(pass string) bool {
	if len(pass) <= 5 {
		return true
	}

	return false
}

func (u *UserRepository) CheckPassAlphabet(pass string) bool {
	for _, charVariable := range pass {
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
			return false
		}
	}
	return true
}
