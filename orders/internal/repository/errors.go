package repository

import "errors"

var (
	ErrIdDoesntExists = errors.New("Entity with Id not found")
)
