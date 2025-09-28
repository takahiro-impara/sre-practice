package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     Email     `json:"email"`
	Password  Password  `json:"password"`
	Name      Name      `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Email string

type Name string

type Password string

// NewUser creates a new user
func NewUser(email Email, password Password, name Name) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate validates the user
func (u *User) Validate() error {
	if err := ValidateEmail(u.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := ValidateName(u.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if err := ValidatePassword(u.Password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

// UpdateName updates the name of the user
func (u *User) UpdateName(name Name) error {
	if err := ValidateName(name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail updates the email of the user
func (u *User) UpdateEmail(email Email) error {
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// UpdatePassword updates the password of the user
func (u *User) UpdatePassword(password Password) error {
	if err := ValidatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	u.Password = password
	u.UpdatedAt = time.Now()
	return nil
}