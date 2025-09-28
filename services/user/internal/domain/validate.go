package domain

import (
	"fmt"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email Email) error {
	if string(email) == "" {
		return fmt.Errorf("email is empty: %w", ErrInvalidEmail)
	}
	if !emailRegex.MatchString(string(email)) {
		return fmt.Errorf("email format is invalid: %w", ErrInvalidEmail)
	}
	return nil
}

func ValidateName(name Name) error {
	if string(name) == "" {
		return fmt.Errorf("name is empty: %w", ErrInvalidName)
	}
	if len(name) < 3 || len(name) > 255 {
		return fmt.Errorf("name length is invalid (len=%d): %w", len(name), ErrInvalidName)
	}
	return nil
}

func ValidatePassword(password Password) error {
	if string(password) == "" {
		return fmt.Errorf("password is empty: %w", ErrInvalidPassword)
	}
	if len(password) < 8 || len(password) > 255 {
		return fmt.Errorf("password length is invalid (len=%d): %w", len(password), ErrInvalidPassword)
	}
	return nil
}
