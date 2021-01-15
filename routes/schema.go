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

type updateAccountSchema struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type makeAdminSchema struct {
	Username string `json:"username" validate:"required"`
}

type revokeAdminURI struct {
	UserID string `uri:"user_id" binding:"required"`
}

type tagSchema struct {
	Name string `json:"name" validate:"required"`
}

// Validate json body
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
