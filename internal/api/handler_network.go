package api

import (
	"net/http"

	"github.com/ptmind/piadmin/internal/network"
)

type networkHandler struct{}

func (h *networkHandler) Interfaces(w http.ResponseWriter, r *http.Request) {
	ifaces, err := network.Interfaces()
	if err != nil {
		http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ifaces)
}
