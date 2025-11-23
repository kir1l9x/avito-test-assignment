package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestPrIDValidate_Direct(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid simple", "pr-1", true},
		{"Valid multi-digit", "pr-12345", true},
		{"Invalid missing prefix", "123", false},
		{"Invalid wrong prefix", "p-123", false},
		{"Invalid uppercase", "PR-1", false},
		{"Invalid no dash", "pr1", false},
		{"Invalid trailing chars", "pr-123a", false},
		{"Empty", "", false},
	}

	v := validator.New()
	_ = RegisterPrIDValidator(v)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.input, "pr_id")
			if tt.expected {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
