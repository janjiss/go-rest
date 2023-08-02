package users

import (
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

type UserRepository interface {
	CreateUser(user *User) error
	GetAllUsers() []User
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repo *UserRepositoryImpl) CreateUser(user *User) error {
	return repo.DB.Create(&user).Error
}

func (repo *UserRepositoryImpl) GetAllUsers() []User {
	var users []User
	repo.DB.Find(&users)

	return users
}
