package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

var validate = validator.New()

// Validate validates the provided data using the validator package
func Validate(data any) error {
	return validate.Struct(data)
}

// BindAndValidate binds and validates the request body against the provided struct
func BindAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return Validate(dst)
}

// RespondWithValidationErrors handles both validation and general errors as JSON
func RespondWithValidationErrors(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if ve, ok := err.(validator.ValidationErrors); ok {
		w.WriteHeader(http.StatusBadRequest)
		var errs []ValidationError
		for _, fe := range ve {
			errs = append(errs, ValidationError{
				Field: fe.Field(),
				Error: fmt.Sprintf("failed on %s", fe.Tag()),
			})
		}
		json.NewEncoder(w).Encode(ValidationErrorResponse{Errors: errs})
		return
	}

	// fallback for non-validation errors (e.g., bad JSON)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
}
