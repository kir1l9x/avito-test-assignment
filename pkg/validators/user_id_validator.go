package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var userIDRegex = regexp.MustCompile(`^u[0-9]+$`)

func RegisterUserIDValidator(v *validator.Validate) error {
	return v.RegisterValidation("user_id", UserIDValidate)
}

func UserIDValidate(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return userIDRegex.MatchString(id)
}
