package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func InitializeInValidator(v *validator.Validate) error {
	if err := RegisterPrIDValidator(v); err != nil {
		return fmt.Errorf("failed to register pr_id validator: %w", err)
	}

	if err := RegisterUserIDValidator(v); err != nil {
		return fmt.Errorf("failed to register user_id validator: %w", err)
	}

	return nil
}

func InitDomainValidator() error {
	Validator = validator.New()
	return nil
}
