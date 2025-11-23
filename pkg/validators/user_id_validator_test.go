package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestUserIDValidate_Direct(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid simple", "u1", true},
		{"Valid multi-digit", "u999", true},
		{"Invalid missing prefix", "1", false},
		{"Invalid wrong prefix", "user-1", false},
		{"Invalid uppercase", "U1", false},
		{"Invalid no digit", "u", false},
		{"Invalid trailing chars", "u123a", false},
		{"Empty", "", false},
	}

	v := validator.New()
	_ = RegisterUserIDValidator(v)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.input, "user_id")
			if tt.expected {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
