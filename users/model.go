package users

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	gormStructParser "janjiss.com/rest/helpers"
	dbHelpers "janjiss.com/rest/helpers/db"
)

var validate *validator.Validate

var DB *gorm.DB

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name  string    `gorm:"type:varchar" json:"name" validate:"required"`
	Email string    `gorm:"type:varchar(100);unique_index" json:"email" validate:"uniquness,required,email"`
}

func newValidator(db *gorm.DB) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("uniquness", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().Interface() // get field value
		fieldKind := fl.Field().Kind()       // get field kind

		newModelStruct := reflect.New(reflect.TypeOf(fl.Top().Interface())).Interface()

		parsedFields, err := gormStructParser.MapStructFieldsToDBFields(newModelStruct)

		if err != nil {
			return false
		}

		// only string fields are supported
		if fieldKind == reflect.String {
			// zero value is not considered a duplicate
			if fieldValue.(string) == "" {
				return true
			}

			// construct a query to count instances of the field value
			query := fmt.Sprintf("%s = ?", parsedFields[fieldName])

			// Create new struct for db.First()
			result := db.Model(fl.Top().Interface()).Where(query, fieldValue).First(newModelStruct)

			return result.Error != nil
		} else {
			return false
		}
	})

	return validate
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {

	validate = newValidator(db)

	err = validate.Struct(user)

	if err != nil {
		return err
	}

	uuid, err := dbHelpers.GenerateULID()

	if err != nil {
		return err
	}

	user.ID = uuid
	return
}
