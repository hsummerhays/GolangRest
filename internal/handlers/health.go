package handlers

import (
	"encoding/json"
	"golangrest/internal/repository"
	"log/slog"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler godoc
// @Summary      Show the status of server
// @Description  get the status of server
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health endpoint accessed", "method", r.Method, "path", r.URL.Path)
	response := HealthResponse{Status: "UP"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ReadyHandler godoc
// @Summary      Show the readiness of server
// @Description  pings the database to verify readiness
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      503  {string}  string "Service Unavailable"
// @Router       /ready [get]
func ReadyHandler(repo repository.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Ready endpoint accessed", "method", r.Method, "path", r.URL.Path)
		if err := repo.Ping(r.Context()); err != nil {
			slog.Error("Readiness check failed", "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		response := HealthResponse{Status: "UP"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			slog.Error("Error encoding response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
