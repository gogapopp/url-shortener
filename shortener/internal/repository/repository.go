package repository

import "errors"

var (
	ErrURLNotExists     = errors.New("url not exists")
	ErrURLAlreadyExists = errors.New("url already exists")
)
