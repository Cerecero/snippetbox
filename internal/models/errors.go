package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("modles: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate emails")
)
