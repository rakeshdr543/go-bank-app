package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUserName     = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidUserFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(val string, minLength, maxLength int) error {
	n := len(val)

	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUserName(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidUserName(val) {
		return fmt.Errorf("must only contain letters, numbers and underscores")
	}

	return nil
}

func ValidateUserFullName(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidUserFullName(val) {
		return fmt.Errorf("must only contain letters and spaces")
	}

	return nil
}

func ValidatePassword(val string) error {
	return ValidateString(val, 6, 40)
}

func ValidateEmail(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(val); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}
