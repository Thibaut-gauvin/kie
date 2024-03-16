package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validateListenPort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		input            string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "valid port 1",
			input:            "8080",
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "default",
			input:            "9145",
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "invalid port: not a number",
			input:            "lorem",
			expectError:      true,
			expectedErrorMsg: "lorem is not a valid port number",
		},
		{
			name:             "invalid port: out of range",
			input:            "66666660",
			expectError:      true,
			expectedErrorMsg: "66666660 is not a valid port number",
		},
		{
			name:             "invalid port: negative out of range",
			input:            "-1",
			expectError:      true,
			expectedErrorMsg: "-1 is not a valid port number",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validateListenPort(test.input)
			if test.expectError {
				assert.EqualError(t, err, test.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
