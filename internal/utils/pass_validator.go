package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const minPasswordLength = 12

// ValidatePasswordWithRegex checks if a password meets complexity requirements using regular expressions.
func ValidatePasswordWithRegex(password string) error {
	// Rule 1: Minimum length check
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}

	// Define regex for each requirement
	rules := map[string]string{
		"an uppercase letter": `[A-Z]`,
		"a lowercase letter": `[a-z]`,
		"a number":            `[0-9]`,
		"a special character": `[!@#$%^&*()_+\-=\[\]{}|;':",.<>/?]`,
	}

	var missing []string
	for requirement, pattern := range rules {
		matched, _ := regexp.MatchString(pattern, password)
		if !matched {
			missing = append(missing, requirement)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("password is missing: %s", strings.Join(missing, ", "))
	}

	return nil
}
