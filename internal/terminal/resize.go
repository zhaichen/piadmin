package terminal

import (
	"encoding/json"
	"os"

	"github.com/creack/pty"
)

type controlMessage struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

func handleControl(ptmx *os.File, msg []byte) {
	var ctrl controlMessage
	if err := json.Unmarshal(msg, &ctrl); err != nil {
		return
	}

	switch ctrl.Type {
	case "resize":
		if ctrl.Cols > 0 && ctrl.Rows > 0 {
			_ = pty.Setsize(ptmx, &pty.Winsize{
				Cols: ctrl.Cols,
				Rows: ctrl.Rows,
			})
		}
	}
}
