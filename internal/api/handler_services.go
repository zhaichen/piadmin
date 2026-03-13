package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ptmind/piadmin/internal/services"
)

type servicesHandler struct{}

func (h *servicesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := services.List()
	if err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, list)
}

func (h *servicesHandler) Status(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	status, err := services.Status(name)
	if err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": status})
}

func (h *servicesHandler) Action(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonErr("invalid request"), http.StatusBadRequest)
		return
	}

	if err := services.Action(name, req.Action); err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}
