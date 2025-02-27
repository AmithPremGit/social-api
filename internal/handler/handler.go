package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"social-api/internal/auth"
	"social-api/internal/cache"
	"social-api/internal/config"
	"social-api/internal/store"
)

// Application contains the application handler dependencies
type Application struct {
	Config        config.Config
	Logger        *log.Logger
	Authenticator *auth.JWTAuthenticator
	Cache         cache.Cache
	UserStore     store.UserStore
	PostStore     store.PostStore
	Validator     *validator.Validate
}

// NewApplication creates a new application handler
func NewApplication(
	cfg config.Config,
	logger *log.Logger,
	authenticator *auth.JWTAuthenticator,
	cache cache.Cache,
	userStore store.UserStore,
	postStore store.PostStore,
) *Application {
	validate := validator.New()

	return &Application{
		Config:        cfg,
		Logger:        logger,
		Authenticator: authenticator,
		Cache:         cache,
		UserStore:     userStore,
		PostStore:     postStore,
		Validator:     validate,
	}
}

// GetIDParam extracts and parses an ID URL parameter
func (app *Application) GetIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id parameter: %s", idParam)
	}
	return id, nil
}

// ValidateRequest validates a request body
func (app *Application) ValidateRequest(v interface{}) error {
	return app.Validator.Struct(v)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// FormatValidationErrors formats validation errors into a readable format
func (app *Application) FormatValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())
			validationErrors = append(validationErrors, ValidationError{
				Field:   field,
				Message: getValidationErrorMessage(e),
			})
		}
	}

	return validationErrors
}

// getValidationErrorMessage returns a user-friendly validation error message
func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", e.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", e.Param())
	default:
		return fmt.Sprintf("Failed validation on %s", e.Tag())
	}
}
