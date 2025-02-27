package handler

import (
	"net/http"

	"social-api/internal/model"
)

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

// HealthCheck handles the health check endpoint
func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Create health check response
	response := HealthCheckResponse{
		Status:  "ok",
		Version: "1.0.0",
		Env:     app.Config.Env,
	}

	// Send response
	err := model.WriteJSON(w, http.StatusOK, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
