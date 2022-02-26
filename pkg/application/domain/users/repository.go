package users

import "errors"

var (
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	GetByUsername(un Username) (*User, error)
}
