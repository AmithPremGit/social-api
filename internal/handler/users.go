package handler

import (
	"errors"
	"net/http"
	"time"

	"social-api/internal/auth"
	"social-api/internal/cache"
	"social-api/internal/model"
	"social-api/internal/store"
)

// RegisterUser handles the user registration endpoint
func (app *Application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input model.UserInput
	err := model.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate input
	err = app.ValidateRequest(input)
	if err != nil {
		validationErrors := app.FormatValidationErrors(err)
		app.validationErrorResponse(w, r, validationErrors)
		return
	}

	// Create user object
	user := &store.User{
		Username: input.Username,
		Email:    input.Email,
		IsActive: true,
	}

	// Set password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create user in database
	err = app.UserStore.Create(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.conflictResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.conflictResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate token
	token, err := app.Authenticator.GenerateToken(user.ID, app.Config.Auth.TokenExpiry)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create response
	response := model.TokenResponse{
		Token: token,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		ExpiresAt: time.Now().Add(app.Config.Auth.TokenExpiry),
	}

	// Cache user if enabled
	if app.Cache != nil {
		err = app.Cache.Set(r.Context(), cache.UserKey(user.ID), user, 1*time.Hour)
		if err != nil {
			app.Logger.Printf("Error caching user: %v", err)
		}
	}

	// Send response
	err = model.WriteJSON(w, http.StatusCreated, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// CreateToken handles the login endpoint
func (app *Application) CreateToken(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input model.UserLoginInput
	err := model.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate input
	err = app.ValidateRequest(input)
	if err != nil {
		validationErrors := app.FormatValidationErrors(err)
		app.validationErrorResponse(w, r, validationErrors)
		return
	}

	// Get user by email
	user, err := app.UserStore.GetByEmail(r.Context(), input.Email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.unauthorizedResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check password
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.unauthorizedResponse(w, r)
		return
	}

	// Generate token
	token, err := app.Authenticator.GenerateToken(user.ID, app.Config.Auth.TokenExpiry)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create response
	response := model.TokenResponse{
		Token: token,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		ExpiresAt: time.Now().Add(app.Config.Auth.TokenExpiry),
	}

	// Cache user if enabled
	if app.Cache != nil {
		err = app.Cache.Set(r.Context(), cache.UserKey(user.ID), user, 1*time.Hour)
		if err != nil {
			app.Logger.Printf("Error caching user: %v", err)
		}
	}

	// Send response
	err = model.WriteJSON(w, http.StatusOK, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetCurrentUser handles the current user endpoint
func (app *Application) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		app.unauthorizedResponse(w, r)
		return
	}

	// Create response
	response := model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// Send response
	err := model.WriteJSON(w, http.StatusOK, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
