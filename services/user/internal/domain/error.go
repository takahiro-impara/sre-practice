package domain

import "fmt"

var (
	ErrInvalidEmail       = NewError("[E001]invalid email")
	ErrInvalidName        = NewError("[E002]invalid name")
	ErrInvalidPassword    = NewError("[E003]invalid password")
	ErrInvalidID          = NewError("[E004]invalid id")
	ErrInvalidLimit       = NewError("[E005]invalid limit")
	ErrInvalidOffset      = NewError("[E006]invalid offset")
	ErrInvalidUpdateInput = NewError("[E007]invalid update input: at least one field must be provided for update")
	ErrNotFound           = NewError("[E008]not found")
	ErrUserNotFound       = NewError("[E009]user not found")
	ErrUserAlreadyExists  = NewError("[E010]user already exists")
	ErrDuplicateID        = NewError("[E011]id already exists")
	ErrDuplicateEmail     = NewError("[E012]email already exists")
	ErrInvalidInput       = NewError("[E013]invalid input")
)

func NewError(message string) error {
	return fmt.Errorf("domain: %s", message)
}
