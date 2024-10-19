package service

import "errors"

var (
	ErrAlreadyExists   = errors.New("user already exists")
	ErrNotFound        = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)
