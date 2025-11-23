package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var prIDRegex = regexp.MustCompile(`^pr-[0-9]+$`)

func RegisterPrIDValidator(v *validator.Validate) error {
	return v.RegisterValidation("pr_id", PrIDValidate)
}

func PrIDValidate(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return prIDRegex.MatchString(id)
}
