package handlers

import (
	"encoding/json"
	"net/http"

	"codeberg.org/hum3/mock-fps/internal/lifecycle"
)

type AdminHandler struct {
	engine *lifecycle.Engine
}

func NewAdminHandler(engine *lifecycle.Engine) *AdminHandler {
	return &AdminHandler{engine: engine}
}

type standInResponse struct {
	Enabled     bool `json:"enabled"`
	QueueLength int  `json:"queue_length"`
}

func (h *AdminHandler) GetStandIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standInResponse{
		Enabled:     h.engine.StandInEnabled(),
		QueueLength: h.engine.QueueLength(),
	})
}

func (h *AdminHandler) SetStandIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	h.engine.SetStandIn(req.Enabled)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standInResponse{
		Enabled:     h.engine.StandInEnabled(),
		QueueLength: h.engine.QueueLength(),
	})
}
