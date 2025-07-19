package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExist    = errors.New("url exist")
)
