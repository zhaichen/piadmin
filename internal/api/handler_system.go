package api

import (
	"encoding/json"
	"net/http"

	"github.com/ptmind/piadmin/internal/monitor"
)

type systemHandler struct {
	collector *monitor.Collector
}

func (h *systemHandler) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	snap := h.collector.GetSnapshot()
	if snap == nil {
		http.Error(w, `{"error":"no data yet"}`, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap)
}
