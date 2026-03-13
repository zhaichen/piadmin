package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ptmind/piadmin/internal/monitor"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type wsHandler struct {
	collector *monitor.Collector
}

func (h *wsHandler) Monitor(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}
	defer conn.Close()

	// send current snapshot immediately
	if snap := h.collector.GetSnapshot(); snap != nil {
		msg, _ := json.Marshal(wsMessage{Type: "snapshot", Data: snap})
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}

	ch := h.collector.Subscribe()
	defer h.collector.Unsubscribe(ch)

	// read pump (detect client disconnect)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	for {
		select {
		case snap, ok := <-ch:
			if !ok {
				return
			}
			msg, _ := json.Marshal(wsMessage{Type: "snapshot", Data: snap})
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}
