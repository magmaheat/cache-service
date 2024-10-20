package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrFileNotFound      = errors.New("file not found")
	ErrInvalidPassword   = errors.New("invalid password")
)
