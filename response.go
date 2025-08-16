package fault

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"net/http"
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

func RepositorError(err error) *APIResponse {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			handleUniqueViolation(c, pqErr)
		case "23503": // foreign_key_violation
			handleForeignKeyViolation(c, pqErr)
		case "23502": // not_null_violation
			handleNotNullViolation(c, pqErr)
		case "23514": // check_violation
			handleCheckViolation(c, pqErr)
		case "42P01": // undefined_table
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "DATABASE_ERROR", // TODO remove these in prod
					Message: "Database table not found",
					Details: map[string]interface{}{
						"error_code": string(pqErr.Code),
					},
				},
			})
		case "42804": // datatype_mismatch
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "INVALID_DATA_TYPE",
					Message: "Invalid data type provided",
					Details: map[string]interface{}{
						"error_code": string(pqErr.Code) + err.Error(),
					},
				},
			})
		case "42703pq":
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "INVALID_FIELD",
					Message: "Invalid field provided",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "DATABASE_ERROR",
					Message: "Database operation failed",
					Details: map[string]interface{}{
						"error_code": string(pqErr.Code) + err.Error(),
					},
				},
			})
		}
		return
	}
}
