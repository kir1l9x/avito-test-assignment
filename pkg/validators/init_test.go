package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestInitializeInValidator_RegistersValidators(t *testing.T) {
	v := validator.New()

	err := InitializeInValidator(v)
	require.NoError(t, err)

	err = v.Var("pr-123", "pr_id")
	require.NoError(t, err)

	err = v.Var("wrong", "pr_id")
	require.Error(t, err)

	err = v.Var("u55", "user_id")
	require.NoError(t, err)

	err = v.Var("wrong", "user_id")
	require.Error(t, err)
}

func TestInitDomainValidator(t *testing.T) {
	err := InitDomainValidator()
	require.NoError(t, err)
	require.NotNil(t, Validator)

	require.Panics(t, func() {
		_ = Validator.Var("pr-1", "pr_id")
	})

	require.Panics(t, func() {
		_ = Validator.Var("u1", "user_id")
	})
}
