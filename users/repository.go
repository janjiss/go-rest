package users

import (
	"fmt"

	"gorm.io/gorm"
	"janjiss.com/rest/helpers/paginator"
	validator "janjiss.com/rest/helpers/validators"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

type UserRepository interface {
	CreateUser(user *User) error
	GetAllUsers(cursor string) ([]User, error)
	FindOneByEmail(email string) (*User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repo *UserRepositoryImpl) CreateUser(user *User) error {
	err := validator.NewValidator(repo.DB).Struct(user)

	if err != nil {
		return err
	}

	return repo.DB.Create(&user).Error
}

func (repo *UserRepositoryImpl) GetAllUsers(cursor string) ([]User, error) {
	users, err := paginator.FetchNextPage[User](repo.DB, cursor, 10)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (repo *UserRepositoryImpl) FindOneByEmail(email string) (*User, error) {
	var user User
	if result := repo.DB.Where("email = ?", email).First(&user); result.Error != nil {
		fmt.Println("Error fetching user:", result.Error)
		return &user, result.Error
	}

	return &user, nil
}
