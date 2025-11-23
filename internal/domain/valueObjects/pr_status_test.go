package valueObjects

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPullRequestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   PullRequestStatus
		expected bool
	}{
		{
			name:     "Valid OPEN",
			status:   OpenPullRequest,
			expected: true,
		},
		{
			name:     "Valid MERGED",
			status:   MergedPullRequest,
			expected: true,
		},
		{
			name:     "Invalid empty",
			status:   PullRequestStatus(""),
			expected: false,
		},
		{
			name:     "Invalid arbitrary",
			status:   PullRequestStatus("SOMETHING"),
			expected: false,
		},
		{
			name:     "Invalid lowercase",
			status:   PullRequestStatus("open"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestPullRequestStatus_String(t *testing.T) {
	require.Equal(t, "OPEN", OpenPullRequest.String())
	require.Equal(t, "MERGED", MergedPullRequest.String())
	require.Equal(t, "INVALID", PullRequestStatus("INVALID").String())
	require.Equal(t, "", PullRequestStatus("").String())
}
