package handlers

import (
	"encoding/json"
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
