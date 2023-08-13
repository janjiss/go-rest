package users

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	auth "janjiss.com/rest/login"
)

type UserServiceImpl struct {
	repo UserRepository
}

type UserService interface {
	CreateUser(name, email string) (*User, error)
	GetAllUsers(cusor string) ([]User, error)
	Login(email string) (string, error)
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

func NewUserService(db *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{
		repo: NewUserRepository(db),
	}
}

func (us *UserServiceImpl) CreateUser(name, email string) (*User, error) {
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

func (us *UserServiceImpl) GetAllUsers(cusor string) ([]User, error) {
	users, err := us.repo.GetAllUsers(cusor)
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

func (us *UserServiceImpl) Login(email string) (string, error) {
	user, err := us.repo.FindOneByEmail(email)

	var token string

	if err != nil {
		return token, errors.New("User not found")
	}

	token, err = auth.GenerateJWT(user.Email)

	if err != nil {
		return token, errors.New("Failed to generate JWT token")
	}

	return token, nil
}
