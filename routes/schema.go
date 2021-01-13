package routes

import (
	"github.com/go-playground/validator/v10"
)

type userSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type userRegisterSchema struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type publicError struct {
	Field   string
	Message string
}

func validateSchema(data interface{}) *[]publicError {
	err := validator.New().Struct(data)
	var errors []publicError

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, publicError{
				Field:   err.Field(),
				Message: err.Error(),
			})
		}
		return &errors
	}

	return nil
}
