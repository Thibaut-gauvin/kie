package logger_test

import (
	"log/slog"
	"testing"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		input         string
		expected      slog.Level
		expectedError string
	}{
		{
			name:          "parse debug",
			input:         "debug",
			expected:      slog.LevelDebug,
			expectedError: "",
		},
		{
			name:          "parse info",
			input:         "info",
			expected:      slog.LevelInfo,
			expectedError: "",
		},
		{
			name:          "parse warn",
			input:         "warn",
			expected:      slog.LevelWarn,
			expectedError: "",
		},
		{
			name:          "parse info",
			input:         "error",
			expected:      slog.LevelError,
			expectedError: "",
		},
		{
			name:          "parse info",
			input:         "info",
			expected:      slog.LevelInfo,
			expectedError: "",
		},
		{
			name:          "wrong log level",
			input:         "lorem",
			expected:      slog.LevelInfo,
			expectedError: "\"lorem\" is not a valid slog level",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual, err := logger.ParseLevel(test.input)
			assert.Equal(t, test.expected, actual)
			if test.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedError)
			}
		})
	}
}
