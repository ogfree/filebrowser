package http

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var shellUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var shellNotAllowed = []byte("Shell not allowed.")

var shellHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	conn, err := shellUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	// Fail fast
	if !d.server.EnableShell || !d.user.Perm.Execute {
		if err := conn.WriteMessage(websocket.TextMessage, shellNotAllowed); err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
		}

		return 0, nil
	}

	cmd := exec.Command(os.Getenv("SHELL"))
	cmd.Env = os.Environ()
	cmd.Dir = d.user.FullPath(r.URL.Path)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer ptmx.Close()

	// Handle pty resizes
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}

			// Read resize message
			_, resizeMessage, err := conn.ReadMessage()
			if err != nil {
				log.Printf("could not read resize message: %v", err)
				continue
			}

			// Parse resize message
			var resize struct {
				Cols uint16 `json:"cols"`
				Rows uint16 `json:"rows"`
			}
			if err := conn.ReadJSON(&resize); err != nil {
				log.Printf("could not parse resize message: %v", err)
				continue
			}

			// Resize pty
			if err := pty.Setsize(ptmx, &pty.Winsize{
				Rows: resize.Rows,
				Cols: resize.Cols,
			}); err != nil {
				log.Printf("could not resize pty: %v", err)
			}
		}
	}()

	// Pipe pty output to websocket
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := ptmx.Read(buffer)
			if err != nil {
				log.Printf("could not read from pty: %v", err)
				conn.Close()
				return
			}

			if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
				log.Printf("could not write to websocket: %v", err)
				return
			}
		}
	}()

	// Pipe websocket input to pty
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("could not read from websocket: %v", err)
			return 0, nil
		}

		if _, err := ptmx.Write(message); err != nil {
			log.Printf("could not write to pty: %v", err)
			return 0, nil
		}
	}
})
