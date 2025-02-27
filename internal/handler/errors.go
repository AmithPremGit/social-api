package handler

import (
	"errors"
	"net/http"

	"social-api/internal/model"
	"social-api/internal/store"
)

// Error response handlers

// respondError is a generic error response helper
func (app *Application) respondError(w http.ResponseWriter, status int, message string) {
	app.Logger.Printf("ERROR: %s (status: %d)", message, status)
	model.WriteJSON(w, status, model.ErrorResponse{Error: message})
}

// notFoundResponse sends a 404 Not Found response
func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.respondError(w, http.StatusNotFound, "The requested resource could not be found")
}

// methodNotAllowedResponse sends a 405 Method Not Allowed response
func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.respondError(w, http.StatusMethodNotAllowed, "The method is not supported for this resource")
}

// badRequestResponse sends a 400 Bad Request response
func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.respondError(w, http.StatusBadRequest, err.Error())
}

// serverErrorResponse sends a 500 Internal Server Error response
func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Printf("INTERNAL SERVER ERROR: %s", err.Error())
	app.respondError(w, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

// validationErrorResponse sends a 422 Unprocessable Entity response for validation errors
func (app *Application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors []ValidationError) {
	app.Logger.Printf("VALIDATION ERROR: %+v", errors)

	// Convert to model validation errors
	validationErrors := make([]model.ValidationError, len(errors))
	for i, err := range errors {
		validationErrors[i] = model.ValidationError{
			Field: err.Field,
			Error: err.Message,
		}
	}

	model.WriteJSON(w, http.StatusUnprocessableEntity, model.ValidationErrorResponse{
		Errors: validationErrors,
	})
}

// unauthorizedResponse sends a 401 Unauthorized response
func (app *Application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	app.respondError(w, http.StatusUnauthorized, "You must be authenticated to access this resource")
}

// forbiddenResponse sends a 403 Forbidden response
func (app *Application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.respondError(w, http.StatusForbidden, "You don't have permission to access this resource")
}

// conflictResponse sends a 409 Conflict response
func (app *Application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.respondError(w, http.StatusConflict, err.Error())
}

// handleError processes common errors and sends the appropriate response
func (app *Application) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, store.ErrDuplicateEmail):
		app.conflictResponse(w, r, err)
	case errors.Is(err, store.ErrDuplicateUsername):
		app.conflictResponse(w, r, err)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
