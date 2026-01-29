package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Validator() *validator.Validate {
	if validate != nil {
		return validate
	}
	validate = validator.New()
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		re := regexp.MustCompile(`^(\+7|8)[0-9]{10}$`)
		return re.MatchString(phone)
	})
	return validate
}
