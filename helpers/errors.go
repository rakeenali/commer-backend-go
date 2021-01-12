package helpers

import "errors"

// const (
// 	ErrNotFound = "models: resource not found"
// )
var (
	ErrNotFound = errors.New("Resource not found")

	ErrInvalidToken = errors.New("Auth Token is not valid")
)
