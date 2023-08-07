package users

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService struct {
	repo UserRepository
}

type CreateUserError struct {
	Errors []struct {
		Field string `json:"field"`
		Error string `json:"error"`
	}
	Message string
}

func (e CreateUserError) Error() string {
	return e.Message
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: NewUserRepository(db),
	}
}

func (us *UserService) CreateUser(name, email string) (*User, error) {
	user := User{
		Name:  name,
		Email: email,
	}

	err := us.repo.CreateUser(&user)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := CreateUserError{}

		for _, fe := range validationErrors {
			errors.Errors = append(errors.Errors, struct {
				Field string `json:"field"`
				Error string `json:"error"`
			}{
				Field: fe.Field(),
				Error: fe.Error(),
			})

		}
		return &user, errors
	}

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (us *UserService) GetAllUsers() []User {
	return us.repo.GetAllUsers()
}
