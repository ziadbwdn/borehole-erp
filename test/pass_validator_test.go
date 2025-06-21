package testing_set

import (
	"boreholedata-ms/internal/utils" // Assuming this path is correct
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestValidatePasswordWithRegex uses a table-driven approach to test various password inputs.
func TestValidatePasswordWithRegex(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		password      string
		expectError   bool
		errorContains []string // Substrings to check for in the error message
	}{
		{
			name:          "Valid Password",
			password:      "ThisIsA_Valid1_Password!",
			expectError:   false,
			errorContains: nil,
		},
		{
			name:          "Valid Password - Minimum Length",
			password:      "ValidPass12!", // 12 characters
			expectError:   false,
			errorContains: nil,
		},
		{
			name:          "Invalid - Too Short",
			password:      "Short1!",
			expectError:   true,
			errorContains: []string{"at least 12 characters long"},
		},
		{
			name:          "Invalid - Missing Uppercase",
			password:      "thisis_invalid1!",
			expectError:   true,
			errorContains: []string{"an uppercase letter"},
		},
		{
			name:          "Invalid - Missing Lowercase",
			password:      "THISIS_INVALID1!",
			expectError:   true,
			errorContains: []string{"a lowercase letter"},
		},
		{
			name:          "Invalid - Missing Number",
			password:      "ThisIs_Invalid!",
			expectError:   true,
			errorContains: []string{"a number"},
		},
		{
			name:          "Invalid - Missing Special Character",
			password:      "ThisIsInvalid1",
			expectError:   true,
			errorContains: []string{"a special character"},
		},
		{
			name:        "Invalid - Missing Multiple Requirements",
			password:    "invalid",
			expectError: true,
			errorContains: []string{
				"at least 12 characters long", // This error is returned first
			},
		},
		{
			name:        "Invalid - Missing Multiple (but long enough)",
			password:    "thispasswordislong",
			expectError: true,
			errorContains: []string{ // These can appear in any order
				"an uppercase letter",
				"a number",
				"a special character",
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the function under test
			err := utils.ValidatePasswordWithRegex(tc.password)

			if tc.expectError {
				// Assert that an error was returned
				assert.Error(t, err, "Expected an error for password: %s", tc.password)
				if err != nil {
					// Assert that the error message contains the expected text
					for _, substring := range tc.errorContains {
						assert.True(t, strings.Contains(err.Error(), substring), "Error message should contain '%s'", substring)
					}
				}
			} else {
				// Assert that no error was returned
				assert.NoError(t, err, "Did not expect an error for password: %s", tc.password)
			}
		})
	}
}