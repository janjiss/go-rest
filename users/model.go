package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	dbHelpers "janjiss.com/rest/helpers/db"
)

var validate *validator.Validate

var DB *gorm.DB

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name  string    `gorm:"type:varchar" json:"name" validate:"required"`
	Email string    `gorm:"type:varchar(100);unique_index" json:"email" validate:"uniquness,required,email"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	uuid, err := dbHelpers.GenerateULID()

	if err != nil {
		return err
	}

	user.ID = uuid
	return
}
