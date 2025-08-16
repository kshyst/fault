package fault

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    int                    `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type validationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message,omitempty"`
}

// ValidatorError sends a validation error response TODO add translation
func ValidatorError(err error) *APIResponse {
	var validationErrors []validationErrorDetail

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			validationErrors = append(validationErrors, validationErrorDetail{
				Field:   fe.Field(),
				Message: fe.Error(),
			})
		}
	}

	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code: 1,
			Details: map[string]interface{}{
				"validation_errors": validationErrors,
			},
		},
	}
}

//func RepositorError(err error) *APIResponse {
//
//}
