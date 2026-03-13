package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

type powerHandler struct{}

func (h *powerHandler) Action(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonErr("invalid request"), http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd
	switch req.Action {
	case "shutdown":
		cmd = exec.Command("sudo", "shutdown", "-h", "now")
	case "reboot":
		cmd = exec.Command("sudo", "reboot")
	default:
		http.Error(w, jsonErr("invalid action: must be 'shutdown' or 'reboot'"), http.StatusBadRequest)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}
