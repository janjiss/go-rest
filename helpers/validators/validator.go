package validator

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	gormStructParser "janjiss.com/rest/helpers"
)

func NewValidator(db *gorm.DB) *validator.Validate {
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
