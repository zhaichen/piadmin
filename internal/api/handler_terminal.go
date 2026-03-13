package api

import (
	"log/slog"
	"net/http"

	"github.com/ptmind/piadmin/internal/terminal"
)

type terminalHandler struct{}

func (h *terminalHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("terminal websocket upgrade failed", "err", err)
		return
	}
	defer conn.Close()

	terminal.HandleWebSocket(conn)
}
