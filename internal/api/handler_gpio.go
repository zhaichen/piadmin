package api

import (
	"encoding/json"
	"net/http"

	"github.com/ptmind/piadmin/internal/gpio"
)

type gpioHandler struct{}

func (h *gpioHandler) Available(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]bool{"available": gpio.Available()})
}

func (h *gpioHandler) List(w http.ResponseWriter, r *http.Request) {
	pins, err := gpio.ListPins()
	if err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, pins)
}

func (h *gpioHandler) SetPin(w http.ResponseWriter, r *http.Request) {
	var req gpio.PinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonErr("invalid request"), http.StatusBadRequest)
		return
	}

	if req.Direction != "" {
		if err := gpio.SetDirection(req.Pin, req.Direction); err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}
	}

	if req.Direction == "out" {
		if err := gpio.SetValue(req.Pin, req.Value); err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *gpioHandler) Export(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pin int `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonErr("invalid request"), http.StatusBadRequest)
		return
	}

	if err := gpio.ExportPin(req.Pin); err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *gpioHandler) Unexport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pin int `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonErr("invalid request"), http.StatusBadRequest)
		return
	}

	if err := gpio.UnexportPin(req.Pin); err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}
