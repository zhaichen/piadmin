package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ptmind/piadmin/internal/process"
)

type processHandler struct{}

func (h *processHandler) List(w http.ResponseWriter, r *http.Request) {
	procs, err := process.List()
	if err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, procs)
}

func (h *processHandler) Kill(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		http.Error(w, jsonErr("invalid pid"), http.StatusBadRequest)
		return
	}
	force := r.URL.Query().Get("force") == "true"

	if err := process.Kill(int32(pid), force); err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonErr(msg string) string {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}
