package helpers

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Date format validation 02/02/1987 date birth
func ValidateDateBirthFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("02/01/2006", fl.Field().String())

	return err == nil
}

func ValidateArrayString(fl validator.FieldLevel) bool {
	field := fl.Field()

	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i).String()
		if err := validator.New().Var(elem, "required,uuid"); err != nil {
			return false
		}
	}

	return true
}
