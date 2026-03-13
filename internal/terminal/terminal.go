package terminal

import (
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

const defaultShell = "/bin/bash"

func getShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	if runtime.GOOS == "darwin" {
		return "/bin/zsh"
	}
	return defaultShell
}

func HandleWebSocket(conn *websocket.Conn) {
	shell := getShell()
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to start shell: "+err.Error()))
		return
	}
	defer ptmx.Close()

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		})
	}
	defer cleanup()

	// pty -> websocket
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				conn.Close()
				return
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// websocket -> pty
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			cleanup()
			return
		}

		switch msgType {
		case websocket.TextMessage:
			// handle resize: JSON message {"type":"resize","cols":80,"rows":24}
			if len(msg) > 0 && msg[0] == '{' {
				handleControl(ptmx, msg)
				continue
			}
			_, _ = ptmx.Write(msg)
		case websocket.BinaryMessage:
			_, _ = ptmx.Write(msg)
		}
	}
}
