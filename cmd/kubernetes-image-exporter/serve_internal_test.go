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

func Test_validateKubeconfigPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		input            string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "empty (default)",
			input:            "",
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "valid file",
			input:            "../../test/fixtures/kubeconfig.yaml",
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "invalid file not exist",
			input:            "/tmp/lorem.txt",
			expectError:      true,
			expectedErrorMsg: "file /tmp/lorem.txt does not exist",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validateKubeconfigPath(test.input)
			if test.expectError {
				assert.EqualError(t, err, test.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_validateRefreshInterval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		input            int
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "default interval",
			input:            30,
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "valid interval",
			input:            15,
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "invalid interval out of range +",
			input:            100000,
			expectError:      true,
			expectedErrorMsg: "100000 is not a valid refresh interval",
		},
		{
			name:             "invalid interval out of range -",
			input:            0,
			expectError:      true,
			expectedErrorMsg: "0 is not a valid refresh interval",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validateRefreshInterval(test.input)
			if test.expectError {
				assert.EqualError(t, err, test.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
